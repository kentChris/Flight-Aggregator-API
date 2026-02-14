package airasia

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

type AirAsiaResponse struct {
	Status  string                 `json:"status"`
	Flights []entity.AirAsiaFlight `json:"flights"`
}

type airAsiaService struct {
	filePath string
}

type AirAsiaService interface {
	GetFlight(ctx context.Context) ([]entity.Flight, error)
}

func NewAirAsiaService(path string) AirAsiaService {
	return &airAsiaService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (a *airAsiaService) GetFlight(ctx context.Context) ([]entity.Flight, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	type result struct {
		flights []entity.Flight
		err     error
	}
	resChan := make(chan result, 1)

	go func() {
		// Mock delay (50 - 150ms)
		delay := time.Duration(rand.Intn(101)+50) * time.Millisecond
		time.Sleep(delay)
		// time.Sleep(4 * time.Second)

		if rand.Intn(100) >= 90 {
			resChan <- result{nil, fmt.Errorf("AirAsia: Random Mock Failure")}
			return
		}

		data, err := os.ReadFile(a.filePath)
		if err != nil {
			resChan <- result{nil, fmt.Errorf("AirAsia.getFlight: Failed to read file")}
			return
		}

		var airAsiaResponse AirAsiaResponse
		if err := json.Unmarshal(data, &airAsiaResponse); err != nil {
			resChan <- result{nil, err}
			return
		}

		flights, err := a.mapFlights(airAsiaResponse.Flights)
		resChan <- result{flights, err}
	}()

	select {
	case res := <-resChan:
		return res.flights, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("AirAsia fetch timed out after 2s")
	}
}

func (a *airAsiaService) mapFlights(rawFlights []entity.AirAsiaFlight) ([]entity.Flight, error) {
	log := logger.Init()
	if len(rawFlights) == 0 {
		return []entity.Flight{}, nil
	}

	unifiedFlights := make([]entity.Flight, 0, len(rawFlights))
	for _, raw := range rawFlights {
		if err := raw.Validate(); err != nil {
			log.Errorf("Validation failed for AirAsia flight: %v", err)
			continue
		}

		unified, err := a.mapFlight(raw)
		if err != nil {
			log.Errorf("Error AirAsia.mapFlights: Corrupted data for flight %s", raw.FlightCode)
			continue
		}
		unifiedFlights = append(unifiedFlights, unified)
	}

	return unifiedFlights, nil
}

func (a *airAsiaService) mapFlight(flight entity.AirAsiaFlight) (entity.Flight, error) {
	depTime, err := time.Parse(time.RFC3339, flight.DepartTime)
	if err != nil {
		return entity.Flight{}, err
	}
	arrTime, err := time.Parse(time.RFC3339, flight.ArriveTime)
	if err != nil {
		return entity.Flight{}, err
	}

	// Convert duration_hours (float64) to minutes
	elapsedDuration := arrTime.Sub(depTime)
	totalMinutes := int(elapsedDuration.Minutes())
	hours := totalMinutes / 60
	mins := totalMinutes % 60
	formattedDuration := fmt.Sprintf("%dh %dm", hours, mins)

	// Initialize Location Registry
	lr := entity.LocationRegistry{}

	// Handle Stops count
	stopCount := len(flight.Stops)

	// Handle Baggage
	baggage := entity.BaggageDetails{
		CarryOn: "No information",
		Checked: "No information",
	}
	parts := strings.Split(flight.BaggageNote, ",")
	if len(parts) >= 2 {
		baggage.CarryOn = strings.TrimSpace(parts[0])
		baggage.Checked = strings.TrimSpace(parts[1])
	} else if len(parts) == 1 {
		baggage.CarryOn = strings.TrimSpace(parts[0])
	}

	return entity.Flight{
		ID:           fmt.Sprintf("%s_%s", flight.FlightCode, entity.AIRASIA),
		Provider:     entity.PROVIDER_AIR_ASIA,
		FlightNumber: flight.FlightCode,
		Airline: entity.AirlineInfo{
			Name: flight.Airline,
			Code: "QZ",
		},
		Departure: entity.LocationDetails{
			Airport:   lr.GetAirport(flight.FromAirport),
			City:      lr.GetCity(flight.FromAirport),
			Datetime:  depTime,
			Timestamp: depTime.Unix(),
			Code:      flight.FromAirport,
		},
		Arrival: entity.LocationDetails{
			Airport:   lr.GetAirport(flight.ToAirport),
			City:      lr.GetCity(flight.ToAirport),
			Datetime:  arrTime,
			Timestamp: arrTime.Unix(),
			Code:      flight.ToAirport,
		},
		Duration: entity.DurationDetails{
			TotalMinutes: totalMinutes,
			Formatted:    formattedDuration,
		},
		Stops: stopCount,
		Price: entity.PriceDetails{
			Amount:    flight.PriceIDR,
			Currency:  "IDR",
			Formatted: util.FormatIDR(flight.PriceIDR),
		},
		AvailableSeats: flight.Seats,
		CabinClass:     flight.CabinClass,
		Aircraft:       nil,
		Baggage:        baggage,
		Amenities:      []string{},
	}, nil
}
