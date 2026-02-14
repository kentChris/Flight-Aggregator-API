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

	return entity.SearchResponse{
		Flights: allFlights,
		SearchCriteria: entity.SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.Date,
			Passengers:    req.Passanger,
			CabinClass:    req.CabinClass,
		},
		Metadata: entity.Metadata{
			TotalResults:       len(allFlights),
			ProvidersQueried:   liveSucceeded + succeeded + failed,
			ProvidersSucceeded: succeeded + liveSucceeded,
			ProvidersFailed:    failed,
			SearchTimeMs:       time.Since(startTime).Milliseconds(),
			CacheHit:           cacheHit,
		},
	}, nil

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
	key := fmt.Sprintf("flights:%s:%s:%s:%s", req.Origin, req.Destination, req.Date, code)

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
	var cachedFlights []entity.Flight
	var missingAirlines []string
	var succeeded int

	for _, code := range airlines {
		var airlineFlights []entity.Flight
		key := fmt.Sprintf("flights:%s:%s:%s:%s", req.Origin, req.Destination, req.Date, code)

		if err := f.redisService.Get(ctx, key, &airlineFlights); err == nil && len(airlineFlights) > 0 {
			cachedFlights = append(cachedFlights, airlineFlights...)
			succeeded++
		} else {
			missingAirlines = append(missingAirlines, code)
		}
	}
	return cachedFlights, missingAirlines, succeeded
}
