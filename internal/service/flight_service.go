package service

import (
	"context"
	logger "flight-aggregator/internal/common"
	"flight-aggregator/internal/entity"
	"flight-aggregator/internal/redis"
	"flight-aggregator/internal/service/airasia"
	"flight-aggregator/internal/service/batikair"
	"flight-aggregator/internal/service/garuda"
	"flight-aggregator/internal/service/lionair"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type flightService struct {
	garudaService   garuda.GarudaService
	batikAirService batikair.BatikAirService
	lionAirService  lionair.LionAirService
	airAsiaService  airasia.AirAsiaService
	redisService    redis.RedisService
}

type FlightService interface {
	SearchFlight(ctx context.Context, req entity.SearchRequest) (entity.SearchResponse, error)
}

func NewFlightService(
	garudaService garuda.GarudaService,
	batikAirService batikair.BatikAirService,
	lionAirService lionair.LionAirService,
	airAsiaService airasia.AirAsiaService,
	redisService redis.RedisService,
) FlightService {
	return &flightService{
		garudaService:   garudaService,
		batikAirService: batikAirService,
		lionAirService:  lionAirService,
		airAsiaService:  airAsiaService,
		redisService:    redisService,
	}
}

func (f *flightService) SearchFlight(ctx context.Context, req entity.SearchRequest) (entity.SearchResponse, error) {
	startTime := time.Now()

	if err := req.Validate(); err != nil {
		return entity.SearchResponse{}, err
	}
	f.standardizeRequest(&req)

	// get from redis if not mark the airlines
	cachedFlights, missingAirlines, succeeded := f.getCachedAirlines(ctx, req)

	fmt.Println(missingAirlines)

	// fetch mock airlines
	var liveFlights []entity.Flight
	var liveSucceeded int
	var failed int
	cacheHit := true

	if len(missingAirlines) != 0 {
		liveFlights, liveSucceeded, failed = f.fetchSpecificAirlines(ctx, req, missingAirlines)
		cacheHit = false
	}
	allFlights := append(cachedFlights, liveFlights...)

	// Fillter
	filteredFlights, bestValue := f.applyFiltersAndIdentifyBest(allFlights, req)

	// Sort
	if len(filteredFlights) > 0 {
		f.applySorting(filteredFlights, req)
	}

	return entity.SearchResponse{
		Flights: filteredFlights,
		SearchCriteria: entity.SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			Passengers:    req.Passanger,
			CabinClass:    req.CabinClass,
		},
		Metadata: entity.Metadata{
			TotalResults:       len(filteredFlights),
			ProvidersQueried:   liveSucceeded + succeeded + failed,
			ProvidersSucceeded: succeeded + liveSucceeded,
			ProvidersFailed:    failed,
			SearchTimeMs:       time.Since(startTime).Milliseconds(),
			CacheHit:           cacheHit,
		},
		BestValue: bestValue,
	}, nil

}

func (f *flightService) applySorting(flights []entity.Flight, req entity.SearchRequest) {
	if len(flights) == 0 {
		return
	}
	sort.Slice(flights, func(i, j int) bool {
		a, b := i, j
		if strings.EqualFold(req.SortOrder, "desc") {
			a, b = j, i
		}

		switch strings.ToLower(req.SortBy) {
		case "duration":
			return flights[a].Duration.TotalMinutes < flights[b].Duration.TotalMinutes
		case "departure_time":
			return flights[a].Departure.Timestamp < flights[b].Departure.Timestamp
		case "arrival_time":
			return flights[a].Arrival.Timestamp < flights[b].Arrival.Timestamp
		default: // Price
			return flights[a].Price.Amount < flights[b].Price.Amount
		}
	})
}

func (f *flightService) applyFiltersAndIdentifyBest(flights []entity.Flight, req entity.SearchRequest) ([]entity.Flight, *entity.Flight) {
	filtered := make([]entity.Flight, 0, len(flights))

	var bestDeal *entity.Flight
	minScore := math.MaxFloat64

	// constant value
	const timeWeight = 2500.0
	const stopPenalty = 150000.0
	const amenities = 50000.0

	for _, fl := range flights {
		if !strings.EqualFold(fl.Departure.Code, req.Origin) ||
			!strings.EqualFold(fl.Arrival.Code, req.Destination) {
			continue
		}

		//  Price, Stop, Duration FILTERS
		if req.PriceMin > 0 && fl.Price.Amount < req.PriceMin {
			continue
		}
		if req.PriceMax > 0 && fl.Price.Amount > req.PriceMax {
			continue
		}
		if req.MaxStops != nil && fl.Stops > *req.MaxStops {
			continue
		}
		if req.MaxDuration > 0 && fl.Duration.TotalMinutes > req.MaxDuration {
			continue
		}

		// Time filters
		depTimeStr := fl.Departure.Datetime.Format("15:04")
		if req.MinDepTime != "" && depTimeStr < req.MinDepTime {
			continue
		}
		if req.MaxDepTime != "" && depTimeStr > req.MaxDepTime {
			continue
		}

		// CALCULATE BEST VALUE SCORE
		// Formula: Price + (Total Time Weight) + (Stop Penalty) - (Amenities)
		score := fl.Price.Amount +
			(float64(fl.Duration.TotalMinutes) * timeWeight) +
			(float64(fl.Stops) * stopPenalty) -
			(float64(len(fl.Amenities)) * amenities)

		if score < minScore {
			minScore = score
			temp := fl
			bestDeal = &temp
		}

		filtered = append(filtered, fl)
	}

	return filtered, bestDeal
}

func (f *flightService) fetchSpecificAirlines(ctx context.Context, req entity.SearchRequest, missingCodes []string) ([]entity.Flight, int, int) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var succeeded, failed int32
	allFlights := []entity.Flight{}
	log := logger.Init()

	providerMap := map[string]func(ctx context.Context) ([]entity.Flight, error){
		entity.GARUDA:   f.garudaService.GetFlight,
		entity.BATIKAIR: f.batikAirService.GetFlight,
		entity.LIONAIR:  f.lionAirService.GetFlight,
		entity.AIRASIA:  f.airAsiaService.GetFlight,
	}

	for _, code := range missingCodes {
		fn, exists := providerMap[code]
		if !exists {
			continue
		}

		wg.Add(1)
		go func(airlineCode string, fetchFn func(context.Context) ([]entity.Flight, error)) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					log.Errorf("Recovered from panic in %s: %v", airlineCode, r)
					atomic.AddInt32(&failed, 1)
				}
			}()

			res, err := fetchFn(ctx)
			if err != nil {
				log.Errorf("API Fetch Failed for %s: %v", airlineCode, err)
				atomic.AddInt32(&failed, 1)
				return
			}

			mu.Lock()
			allFlights = append(allFlights, res...)
			mu.Unlock()
			atomic.AddInt32(&succeeded, 1)

			f.saveToCache(context.Background(), req, airlineCode, res)

		}(code, fn)
	}

	wg.Wait()
	return allFlights, int(succeeded), int(failed)
}

func (f *flightService) saveToCache(ctx context.Context, req entity.SearchRequest, code string, flights []entity.Flight) {
	key := fmt.Sprintf("flights:%s:%s:%s", req.Origin, req.DepartureDate, code)

	err := f.redisService.Set(ctx, key, flights, 1*time.Minute)
	if err != nil {
		logger.Init().Errorf("Redis Save Failed for %s: %v", code, err)
	}
}

func (f *flightService) standardizeRequest(req *entity.SearchRequest) {
	req.Origin = strings.ToUpper(req.Origin)
	req.Destination = strings.ToUpper(req.Destination)
}

func (f *flightService) getCachedAirlines(ctx context.Context, req entity.SearchRequest) ([]entity.Flight, []string, int) {
	airlines := []string{entity.GARUDA, entity.LIONAIR, entity.BATIKAIR, entity.AIRASIA}
	targetAirlines := airlines
	var cachedFlights []entity.Flight
	var missingAirlines []string
	var succeeded int

	if len(req.Airlines) > 0 {
		targetAirlines = req.Airlines
	}

	for _, code := range targetAirlines {
		var airlineFlights []entity.Flight
		key := fmt.Sprintf("flights:%s:%s:%s", req.Origin, req.DepartureDate, code)

		if err := f.redisService.Get(ctx, key, &airlineFlights); err == nil && len(airlineFlights) > 0 {
			cachedFlights = append(cachedFlights, airlineFlights...)
			succeeded++
		} else {
			missingAirlines = append(missingAirlines, code)
		}
	}
	return cachedFlights, missingAirlines, succeeded
}
