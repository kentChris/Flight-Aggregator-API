package entity

import "time"

const GARUDA = "Garuda"

type Flight struct {
	ID             string          `json:"id"`
	Provider       string          `json:"provider"`
	Airline        AirlineInfo     `json:"airline"`
	FlightNumber   string          `json:"flight_number"`
	Departure      LocationDetails `json:"departure"`
	Arrival        LocationDetails `json:"arrival"`
	Duration       DurationDetails `json:"duration"`
	Stops          int             `json:"stops"`
	Price          PriceDetails    `json:"price"`
	AvailableSeats int             `json:"available_seats"`
	CabinClass     string          `json:"cabin_class"`
	Aircraft       *string         `json:"aircraft"`
	Amenities      []string        `json:"amenities"`
	Baggage        BaggageDetails  `json:"baggage"`
}

type AirlineInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type LocationDetails struct {
	Airport   string    `json:"airport"`
	City      string    `json:"city"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

type DurationDetails struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type PriceDetails struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type BaggageDetails struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
