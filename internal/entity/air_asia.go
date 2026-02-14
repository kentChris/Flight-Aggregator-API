package entity

import (
	"fmt"
	"time"
)

type AirAsiaFlight struct {
	FlightCode    string        `json:"flight_code"`
	Airline       string        `json:"airline"`
	FromAirport   string        `json:"from_airport"`
	ToAirport     string        `json:"to_airport"`
	DepartTime    string        `json:"depart_time"`
	ArriveTime    string        `json:"arrive_time"`
	DurationHours float64       `json:"duration_hours"`
	DirectFlight  bool          `json:"direct_flight"`
	Stops         []AirAsiaStop `json:"stops,omitempty"`
	PriceIDR      float64       `json:"price_idr"`
	Seats         int           `json:"seats"`
	CabinClass    string        `json:"cabin_class"`
	BaggageNote   string        `json:"baggage_note"`
}

type AirAsiaStop struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}

func (f *AirAsiaFlight) Validate() error {
	depTime, err := time.Parse(time.RFC3339, f.DepartTime)
	if err != nil {
		return fmt.Errorf("invalid departure time format: %w", err)
	}
	arrTime, err := time.Parse(time.RFC3339, f.ArriveTime)
	if err != nil {
		return fmt.Errorf("invalid arrival time format: %w", err)
	}

	// arrival must be after the departure
	if !arrTime.After(depTime) {
		return fmt.Errorf("flight %s arrives (%v) before it departs (%v)", f.FlightCode, arrTime, depTime)
	}

	// Duration must be positive
	if f.DurationHours <= 0 {
		return fmt.Errorf("flight %s has invalid duration: %.2f hours", f.FlightCode, f.DurationHours)
	}

	// Origin and Destination cannot be the same
	if f.FromAirport == f.ToAirport {
		return fmt.Errorf("flight %s origin and destination are the same: %s", f.FlightCode, f.FromAirport)
	}

	// Price and Seats
	if f.PriceIDR < 0 {
		return fmt.Errorf("flight %s has negative price", f.FlightCode)
	}

	return nil
}
