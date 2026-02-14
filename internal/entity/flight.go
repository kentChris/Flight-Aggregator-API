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
	Code      string
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
	Origin        string   `json:"origin"`
	Destination   []string `json:"destinations"`
	DepartureDate string   `json:"departureDate"`
	ReturnDate    *string  `json:"returnDate"` // Pointer because it can be null
	Passanger     int      `json:"passengers"`
	CabinClass    string   `json:"cabinClass"`

	PriceMin    float64  `json:"priceMin,omitempty"`
	PriceMax    float64  `json:"priceMax,omitempty"`
	MaxStops    *int     `json:"maxStops,omitempty"`
	Airlines    []string `json:"airlines,omitempty"`
	MinDepTime  string   `json:"minDepTime,omitempty"`
	MaxDepTime  string   `json:"maxDepTime,omitempty"`
	MaxDuration int      `json:"maxDuration,omitempty"`

	// Sorting
	SortBy    string `json:"sortBy,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
}

func (r *SearchRequest) Validate() error {
	if len(r.Origin) != 3 {
		return fmt.Errorf("origin must be a 3-letter IATA code")
	}

	if len(r.Destination) == 0 {
		return fmt.Errorf("at least one destination must be provided")
	}

	for _, dest := range r.Destination {
		if len(dest) != 3 {
			return fmt.Errorf("destination %s must be a 3-letter IATA code", dest)
		}

		if r.Origin == dest {
			return fmt.Errorf("origin and destination %s cannot be the same", dest)
		}
	}

	return nil
}
