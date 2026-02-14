package lionair

import (
	"context"
	"encoding/json"
	logger "flight-aggregator/internal/common"
	"flight-aggregator/internal/entity"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type LionResponse struct {
	Success bool     `json:"success"`
	Data    LionData `json:"data"`
}

type LionData struct {
	AvailableFlights []entity.LionFlight `json:"available_flights"`
}

type lionAirService struct {
	filePath string
}

type LionAirService interface {
	GetFlight(ctx context.Context) ([]entity.Flight, error)
}

func NewLionAirService(path string) LionAirService {
	return &lionAirService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (g *lionAirService) GetFlight(ctx context.Context) ([]entity.Flight, error) {
	// ToDo: add utils function
	//mock sleep 100 - 200 ms
	delay := time.Duration(rand.Intn(101)+100) * time.Millisecond
	time.Sleep(delay)

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return nil, fmt.Errorf("Lionair.getFlight: Failed to get Lion air data")
	}

	var response LionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return g.mapFlights(response.Data.AvailableFlights)
}

func (s *lionAirService) mapFlights(rawFlights []entity.LionFlight) ([]entity.Flight, error) {
	log := logger.Init()
	if len(rawFlights) == 0 {
		return []entity.Flight{}, nil
	}

	unifiedFlights := make([]entity.Flight, 0, len(rawFlights))
	for _, raw := range rawFlights {
		if err := raw.Validate(); err != nil {
			log.Errorf("LionAir Integrity Error: %v", err)
			continue
		}

		unified, err := s.mapFlight(raw)
		if err != nil {
			log.Errorf("Error Lion.mapFlights: Corrupted data")
			continue
		}
		unifiedFlights = append(unifiedFlights, unified)
	}

	return unifiedFlights, nil
}

func (s *lionAirService) mapFlight(flight entity.LionFlight) (entity.Flight, error) {
	// Departure
	locDep, _ := time.LoadLocation(flight.Schedule.DepartureTimezone)
	depTime, err := time.ParseInLocation("2006-01-02T15:04:05", flight.Schedule.Departure, locDep)
	if err != nil {
		return entity.Flight{}, err
	}

	// Arrival
	locArr, _ := time.LoadLocation(flight.Schedule.ArrivalTimezone)
	arrTime, err := time.ParseInLocation("2006-01-02T15:04:05", flight.Schedule.Arrival, locArr)
	if err != nil {
		return entity.Flight{}, err
	}

	elapsedDuration := arrTime.Sub(depTime)
	totalMinutes := int(elapsedDuration.Minutes())
	hours := totalMinutes / 60
	mins := totalMinutes % 60
	formattedDuration := fmt.Sprintf("%dh %dm", hours, mins)

	aircraft := flight.PlaneType

	// amenities
	amenities := []string{}
	if flight.Services.WifiAvailable {
		amenities = append(amenities, entity.AMENITIES_WIFI)
	}
	if flight.Services.MealsIncluded {
		amenities = append(amenities, entity.AMENITIES_MEAL)
	}

	return entity.Flight{
		ID:       fmt.Sprintf("%s_%s", flight.ID, entity.LIONAIR),
		Provider: entity.PROVIDER_LION_AIR,
		Airline: entity.AirlineInfo{
			Name: flight.Carrier.Name,
			Code: flight.ID,
		},
		FlightNumber: flight.ID,
		Departure: entity.LocationDetails{
			Airport:   flight.Route.From.Name,
			City:      flight.Route.From.City,
			Datetime:  depTime,
			Timestamp: depTime.Unix(),
			Code:      flight.Route.From.Code,
		},
		Arrival: entity.LocationDetails{
			Airport:   flight.Route.To.Name,
			City:      flight.Route.To.City,
			Datetime:  arrTime,
			Timestamp: arrTime.Unix(),
			Code:      flight.Route.To.Code,
		},
		Duration: entity.DurationDetails{
			TotalMinutes: totalMinutes,
			Formatted:    formattedDuration,
		},
		Stops:          flight.StopCount,
		AvailableSeats: flight.SeatsLeft,
		CabinClass:     flight.Pricing.FareType,
		Aircraft:       &aircraft,
		Price: entity.PriceDetails{
			Amount:   flight.Pricing.Total,
			Currency: flight.Pricing.Currency,
		},
		Baggage: entity.BaggageDetails{
			CarryOn: flight.Services.BaggageAllowance.Cabin,
			Checked: flight.Services.BaggageAllowance.Hold,
		},
		Amenities: amenities,
	}, nil
}
