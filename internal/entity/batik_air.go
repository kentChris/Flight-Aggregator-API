package entity

import (
	"fmt"
	"time"
)

type BatikFlight struct {
	FlightNumber      string            `json:"flightNumber"`
	AirlineName       string            `json:"airlineName"`
	AirlineIATA       string            `json:"airlineIATA"`
	Origin            string            `json:"origin"`
	Destination       string            `json:"destination"`
	DepartureDateTime string            `json:"departureDateTime"`
	ArrivalDateTime   string            `json:"arrivalDateTime"`
	TravelTime        string            `json:"travelTime"`
	NumberOfStops     int               `json:"numberOfStops"`
	Connections       []BatikConnection `json:"connections,omitempty"`
	Fare              BatikFare         `json:"fare"`
	SeatsAvailable    int               `json:"seatsAvailable"`
	AircraftModel     string            `json:"aircraftModel"`
	BaggageInfo       string            `json:"baggageInfo"`
	OnboardServices   []string          `json:"onboardServices"`
}

type BatikFare struct {
	BasePrice    float64 `json:"basePrice"`
	Taxes        float64 `json:"taxes"`
	TotalPrice   float64 `json:"totalPrice"`
	CurrencyCode string  `json:"currencyCode"`
	Class        string  `json:"class"`
}

type BatikConnection struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}

func (f *BatikFlight) Validate() error {
	const layout = "2006-01-02T15:04:05-0700"
	depTime, err := time.Parse(layout, f.DepartureDateTime)
	if err != nil {
		return fmt.Errorf("invalid Batik departure format: %w", err)
	}
	arrTime, err := time.Parse(layout, f.ArrivalDateTime)
	if err != nil {
		return fmt.Errorf("invalid Batik arrival format: %w", err)
	}

	if !arrTime.After(depTime) {
		return fmt.Errorf("[%s] Arrival must be after departure", f.FlightNumber)
	}

	if f.Fare.TotalPrice <= 0 {
		return fmt.Errorf("[%s] Total price must be greater than zero", f.FlightNumber)
	}

	if f.Origin == f.Destination {
		return fmt.Errorf("[%s] Origin and Destination cannot be identical (%s)", f.FlightNumber, f.Origin)
	}

	if f.SeatsAvailable < 0 {
		return fmt.Errorf("[%s] Invalid seat count: %d", f.FlightNumber, f.SeatsAvailable)
	}

	return nil
}
