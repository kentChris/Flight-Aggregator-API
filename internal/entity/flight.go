package entity

import (
	"fmt"
	"time"
)

const GARUDA = "Garuda"
const LIONAIR = "LionAir"
const BATIKAIR = "BatikAir"
const AIRASIA = "AirAsia"

const PROVIDER_LION_AIR = "Lion Air"
const PROVIDER_GARUDA = "Garuda Indonesia"
const PROVIDER_BATIK_AIR = "Batik Air"
const PROVIDER_AIR_ASIA = "Air ASIA"

const AMENITIES_WIFI = "wifi"
const AMENITIES_POWER_OUTLET = "power_outlet"
const AMENITIES_MEAL = "meal"
const AMENITIES_ENTERTAINMENT = "entertainment"

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

type LocationRegistry struct{}

var airportData = map[string]struct {
	City    string
	Airport string
}{
	"CGK": {City: "Jakarta", Airport: "Soekarno-Hatta International"},
	"DPS": {City: "Denpasar", Airport: "Ngurah Rai International"},
	"SUB": {City: "Surabaya", Airport: "Juanda International"},
	"UPG": {City: "Makassar", Airport: "Sultan Hasanuddin International"},
	"SOC": {City: "Solo", Airport: "Adi Soemarmo International"},
}

func (r *LocationRegistry) GetCity(code string) string {
	if loc, ok := airportData[code]; ok {
		return loc.City
	}
	return ""
}

func (r *LocationRegistry) GetAirport(code string) string {
	if loc, ok := airportData[code]; ok {
		return loc.Airport
	}
	return ""
}

type SearchRequest struct {
	Origin      string
	Destination string
	Date        string // Expecting YYYY-MM-DD
	CabinClass  string
	Passanger   int
	PriceMax    float64
	SortBy      string
}

func (r *SearchRequest) Validate() error {
	if len(r.Origin) != 3 {
		return fmt.Errorf("origin must be a 3-letter IATA code")
	}

	if len(r.Destination) != 3 {
		return fmt.Errorf("destination must be a 3-letter IATA code")
	}

	if r.Origin == r.Destination {
		return fmt.Errorf("origin and destination cannot be the same")
	}

	// Commented because the mock data is 2025-12-15
	// searchDate, err := time.Parse("2006-01-02", r.Date)
	// if err != nil {
	// 	return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	// }

	// today := time.Now().Truncate(24 * time.Hour)
	// if searchDate.Before(today) {
	// 	return fmt.Errorf("cannot search for flights in the past")
	// }

	return nil
}
