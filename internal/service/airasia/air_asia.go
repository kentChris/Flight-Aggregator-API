package airasia

import (
	"context"
	"encoding/json"
	logger "flight-aggregator/internal/common"
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
	//mock sleep 50 - 150 ms
	// ToDo: add utils function
	delay := time.Duration(rand.Intn(101)+50) * time.Millisecond
	time.Sleep(delay)

	// ToDo: add this in util
	change := rand.Intn(100)
	if change >= 90 {
		return nil, fmt.Errorf("Mock Fail")
	}

	data, err := os.ReadFile(a.filePath)
	if err != nil {
		return nil, fmt.Errorf("airAsia.getFlight: Failed to get Air Asia data")
	}

	var airAsiaResponse AirAsiaResponse
	if err := json.Unmarshal(data, &airAsiaResponse); err != nil {
		return nil, err
	}

	return a.mapFlights(airAsiaResponse.Flights)
}

func (a *airAsiaService) mapFlights(rawFlights []entity.AirAsiaFlight) ([]entity.Flight, error) {
	log := logger.Init()
	if len(rawFlights) == 0 {
		return []entity.Flight{}, nil
	}

	unifiedFlights := make([]entity.Flight, 0, len(rawFlights))
	for _, raw := range rawFlights {
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
		},
		Arrival: entity.LocationDetails{
			Airport:   lr.GetAirport(flight.ToAirport),
			City:      lr.GetCity(flight.ToAirport),
			Datetime:  arrTime,
			Timestamp: arrTime.Unix(),
		},
		Duration: entity.DurationDetails{
			TotalMinutes: totalMinutes,
			Formatted:    formattedDuration,
		},
		Stops: stopCount,
		Price: entity.PriceDetails{
			Amount:   flight.PriceIDR,
			Currency: "IDR",
		},
		AvailableSeats: flight.Seats,
		CabinClass:     flight.CabinClass,
		Aircraft:       nil,
		Baggage:        baggage,
		Amenities:      []string{},
	}, nil
}
