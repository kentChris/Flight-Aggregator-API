package entity

import (
	"fmt"
	"time"
)

type Departure struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type Arrival struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Baggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type GarudaFlight struct {
	FlightID        string    `json:"flight_id"`
	Airline         string    `json:"airline"`
	AirlineCode     string    `json:"airline_code"`
	Departure       Departure `json:"departure"`
	Arrival         Arrival   `json:"arrival"`
	DurationMinutes int       `json:"duration_minutes"`
	Stops           int       `json:"stops"`
	Aircraft        string    `json:"aircraft"`
	Price           Price     `json:"price"`
	AvailableSeats  int       `json:"available_seats"`
	FareClass       string    `json:"fare_class"`
	Baggage         Baggage   `json:"baggage"`
	Amenities       []string  `json:"amenities"`
}

func (f *GarudaFlight) Validate() error {
	depTime, err := time.Parse(time.RFC3339, f.Departure.Time)
	if err != nil {
		return fmt.Errorf("[%s] invalid departure time: %w", f.FlightID, err)
	}
	arrTime, err := time.Parse(time.RFC3339, f.Arrival.Time)
	if err != nil {
		return fmt.Errorf("[%s] invalid arrival time: %w", f.FlightID, err)
	}

	if !arrTime.After(depTime) {
		return fmt.Errorf("[%s] arrival time (%v) must be after departure time (%v)",
			f.FlightID, arrTime, depTime)
	}

	if f.Departure.Airport == f.Arrival.Airport {
		return fmt.Errorf("[%s] origin and destination airports are the same: %s",
			f.FlightID, f.Departure.Airport)
	}

	if f.Price.Amount <= 0 {
		return fmt.Errorf("[%s] price must be greater than zero", f.FlightID)
	}
	if f.AvailableSeats < 0 {
		return fmt.Errorf("[%s] available seats cannot be negative", f.FlightID)
	}

	return nil
}
