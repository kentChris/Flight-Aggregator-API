package garuda

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

type GarudaResponse struct {
	Status  string                `json:"status"`
	Flights []entity.GarudaFlight `json:"flights"`
}

type garudaService struct {
	filePath string
}

type GarudaService interface {
	GetFlight(ctx context.Context) ([]entity.Flight, error)
}

func NewGarudaService(path string) GarudaService {
	return &garudaService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (g *garudaService) GetFlight(ctx context.Context) ([]entity.Flight, error) {
	//mock sleep 50 - 100 ms
	delay := time.Duration(rand.Intn(51)+50) * time.Millisecond
	// time.Sleep(20 * time.Second)
	time.Sleep(delay)

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return nil, fmt.Errorf("Garuda.getFlight: Failed to get garuda data")
	}

	var garudaResponse GarudaResponse
	if err := json.Unmarshal(data, &garudaResponse); err != nil {
		return nil, err
	}

	return g.mapFlights(garudaResponse.Flights)
}

func (g *garudaService) mapFlights(rawFlights []entity.GarudaFlight) ([]entity.Flight, error) {
	log := logger.Init()
	if len(rawFlights) == 0 {
		return []entity.Flight{}, nil
	}

	unifiedFlights := make([]entity.Flight, 0, len(rawFlights))

	for _, raw := range rawFlights {
		unified, err := g.mapFlight(raw)
		if err != nil {
			log.Errorf("Error Garuda.mapFlights: Corrupted data")
			continue
		}
		unifiedFlights = append(unifiedFlights, unified)
	}

	return unifiedFlights, nil
}

func (g *garudaService) mapFlight(flight entity.GarudaFlight) (entity.Flight, error) {
	// ToDo: handle WIT, WIB, WITA
	depTime, err := time.Parse(time.RFC3339, flight.Departure.Time)
	if err != nil {
		return entity.Flight{}, err
	}
	arrTime, err := time.Parse(time.RFC3339, flight.Arrival.Time)
	if err != nil {
		return entity.Flight{}, err
	}

	elapsedDuration := arrTime.Sub(depTime)
	totalMinutes := int(elapsedDuration.Minutes())
	hours := totalMinutes / 60
	mins := totalMinutes % 60
	formattedDuration := fmt.Sprintf("%dh %dm", hours, mins)

	// init location registery
	locationRegistery := entity.LocationRegistry{}

	return entity.Flight{
		ID:       fmt.Sprintf("%s_%s", flight.FlightID, entity.GARUDA),
		Provider: entity.PROVIDER_GARUDA,
		Airline: entity.AirlineInfo{
			Name: flight.Airline,
			Code: flight.AirlineCode,
		},
		FlightNumber: flight.FlightID,
		Departure: entity.LocationDetails{
			Airport:   locationRegistery.GetAirport(flight.Departure.Airport),
			City:      flight.Departure.City,
			Datetime:  depTime,
			Timestamp: depTime.Unix(),
			Code:      flight.Departure.Airport,
		},
		Arrival: entity.LocationDetails{
			Airport:   locationRegistery.GetAirport(flight.Arrival.Airport),
			City:      flight.Arrival.City,
			Datetime:  arrTime,
			Timestamp: arrTime.Unix(),
			Code:      flight.Arrival.Airport,
		},
		Duration: entity.DurationDetails{
			TotalMinutes: totalMinutes,
			Formatted:    formattedDuration,
		},
		Stops: flight.Stops,
		Price: entity.PriceDetails{
			Amount:   flight.Price.Amount,
			Currency: flight.Price.Currency,
		},
		AvailableSeats: flight.AvailableSeats,
		CabinClass:     flight.FareClass,
		Aircraft:       &flight.Aircraft,
		Amenities:      flight.Amenities,
		Baggage: entity.BaggageDetails{
			CarryOn: fmt.Sprintf("%d piece(s)", flight.Baggage.CarryOn),
			Checked: fmt.Sprintf("%d piece(s)", flight.Baggage.Checked),
		},
	}, nil
}
