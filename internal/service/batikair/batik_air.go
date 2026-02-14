package batikair

import (
	"context"
	"encoding/json"
	logger "flight-aggregator/internal/common"
	"flight-aggregator/internal/common/util"
	"flight-aggregator/internal/entity"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type BatikAirResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Results []entity.BatikFlight `json:"results"`
}

type batikAirService struct {
	filePath string
}

type BatikAirService interface {
	GetFlight(ctx context.Context) ([]entity.Flight, error)
}

func NewBatikAirService(path string) BatikAirService {
	return &batikAirService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (b *batikAirService) GetFlight(ctx context.Context) ([]entity.Flight, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	type result struct {
		flights []entity.Flight
		err     error
	}
	resChan := make(chan result, 1)

	go func() {
		// Mock delay logic (200 - 400ms)
		delay := time.Duration(rand.Intn(201)+200) * time.Millisecond
		time.Sleep(delay)
		// time.Sleep(4 * time.Second)

		data, err := os.ReadFile(b.filePath)
		if err != nil {
			resChan <- result{nil, fmt.Errorf("Batik.getFlight: Failed to read file")}
			return
		}

		var batikAirResponse BatikAirResponse
		if err := json.Unmarshal(data, &batikAirResponse); err != nil {
			resChan <- result{nil, err}
			return
		}

		flights, err := b.mapFlights(batikAirResponse.Results)
		resChan <- result{flights, err}
	}()

	select {
	case res := <-resChan:
		return res.flights, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("BatikAir fetch timed out after 2s")
	}
}

func (b *batikAirService) mapFlights(rawFlights []entity.BatikFlight) ([]entity.Flight, error) {
	log := logger.Init()
	if len(rawFlights) == 0 {
		return []entity.Flight{}, nil
	}
	unifiedFlights := make([]entity.Flight, 0, len(rawFlights))

	for _, raw := range rawFlights {
		if err := raw.Validate(); err != nil {
			log.Errorf("Batik Data Integrity Error: %v", err)
			continue
		}

		unified, err := b.mapFlight(raw)
		if err != nil {
			log.Errorf("Error Lion.mapFlights: Corrupted data")
			continue
		}
		unifiedFlights = append(unifiedFlights, unified)
	}

	return unifiedFlights, nil
}

func (b *batikAirService) mapFlight(flight entity.BatikFlight) (entity.Flight, error) {
	const layout = "2006-01-02T15:04:05-0700"
	depTime, err := time.Parse(layout, flight.DepartureDateTime)
	if err != nil {
		return entity.Flight{}, err
	}
	arrTime, err := time.Parse(layout, flight.ArrivalDateTime)
	if err != nil {
		return entity.Flight{}, err
	}

	elapsedDuration := arrTime.Sub(depTime)
	totalMinutes := int(elapsedDuration.Minutes())
	hours := totalMinutes / 60
	mins := totalMinutes % 60
	formattedDuration := fmt.Sprintf("%dh %dm", hours, mins)

	baggage := entity.BaggageDetails{
		CarryOn: "7kg",
		Checked: "20kg",
	}
	if parts := strings.Split(flight.BaggageInfo, ","); len(parts) == 2 {
		baggage.CarryOn = strings.TrimSpace(strings.ReplaceAll(parts[0], "cabin", ""))
		baggage.Checked = strings.TrimSpace(strings.ReplaceAll(parts[1], "checked", ""))
	}

	aircraft := flight.AircraftModel

	// init location registery
	locationRegistery := entity.LocationRegistry{}

	// formated IDR
	const idr = "IDR"
	var formattedPrice string
	if flight.Fare.CurrencyCode == idr {
		formattedPrice = util.FormatIDR(float64(flight.Fare.TotalPrice))
	}

	return entity.Flight{
		ID:       fmt.Sprintf("%s_%s", flight.FlightNumber, entity.BATIKAIR),
		Provider: entity.PROVIDER_BATIK_AIR,
		Airline: entity.AirlineInfo{
			Name: flight.AirlineName,
			Code: flight.AirlineIATA,
		},
		FlightNumber: flight.FlightNumber,
		Departure: entity.LocationDetails{
			Airport:   locationRegistery.GetAirport(flight.Origin),
			City:      locationRegistery.GetCity(flight.Origin),
			Datetime:  depTime,
			Timestamp: depTime.Unix(),
			Code:      flight.Origin,
		},
		Arrival: entity.LocationDetails{
			Airport:   locationRegistery.GetAirport(flight.Destination),
			City:      locationRegistery.GetCity(flight.Destination),
			Datetime:  arrTime,
			Timestamp: arrTime.Unix(),
			Code:      flight.Destination,
		},
		Duration: entity.DurationDetails{
			TotalMinutes: totalMinutes,
			Formatted:    formattedDuration,
		},
		Stops:          flight.NumberOfStops,
		AvailableSeats: flight.SeatsAvailable,
		CabinClass:     flight.Fare.Class,
		Aircraft:       &aircraft,
		Price: entity.PriceDetails{
			Amount:    flight.Fare.TotalPrice,
			Currency:  flight.Fare.CurrencyCode,
			Formatted: formattedPrice,
		},
		Baggage:   baggage,
		Amenities: flight.OnboardServices,
	}, nil
}
