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
	// 1. Set the 2-second deadline
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	type result struct {
		flights []entity.Flight
		err     error
	}
	resChan := make(chan result, 1)

	go func() {
		// Mock delay logic
		delay := time.Duration(rand.Intn(51)+50) * time.Millisecond
		time.Sleep(delay)
		// time.Sleep(5 * time.Second)

		data, err := os.ReadFile(g.filePath)
		if err != nil {
			resChan <- result{nil, fmt.Errorf("Garuda.getFlight: Failed to read file")}
			return
		}

		var garudaResponse GarudaResponse
		if err := json.Unmarshal(data, &garudaResponse); err != nil {
			resChan <- result{nil, err}
			return
		}

		flights, err := g.mapFlights(garudaResponse.Flights)
		resChan <- result{flights, err}
	}()

	select {
	case res := <-resChan:
		return res.flights, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("Garuda fetch timed out after 2s")
	}
}

func (g *garudaService) mapFlights(rawFlights []entity.GarudaFlight) ([]entity.Flight, error) {
	log := logger.Init()
	if len(rawFlights) == 0 {
		return []entity.Flight{}, nil
	}

	unifiedFlights := make([]entity.Flight, 0, len(rawFlights))

	for _, raw := range rawFlights {

		if err := raw.Validate(); err != nil {
			log.Errorf("Garuda Data Integrity Error: %v", err)
			continue
		}

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
