package entity

import (
	"fmt"
	"time"
)

type LionFlight struct {
	ID         string       `json:"id"`
	Carrier    LionCarrier  `json:"carrier"`
	Route      LionRoute    `json:"route"`
	Schedule   LionSchedule `json:"schedule"`
	FlightTime int          `json:"flight_time"`
	IsDirect   bool         `json:"is_direct"`
	StopCount  int          `json:"stop_count"`
	Layovers   []LionStop   `json:"layovers,omitempty"`
	Pricing    LionPricing  `json:"pricing"`
	SeatsLeft  int          `json:"seats_left"`
	PlaneType  string       `json:"plane_type"`
	Services   LionServices `json:"services"`
}

type LionCarrier struct {
	Name string `json:"name"`
	IATA string `json:"iata"`
}

type LionRoute struct {
	From LionLocation `json:"from"`
	To   LionLocation `json:"to"`
}

type LionLocation struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type LionSchedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type LionPricing struct {
	Total    float64 `json:"total"`
	Currency string  `json:"currency"`
	FareType string  `json:"fare_type"`
}

type LionServices struct {
	WifiAvailable    bool        `json:"wifi_available"`
	MealsIncluded    bool        `json:"meals_included"`
	BaggageAllowance LionBaggage `json:"baggage_allowance"`
}

type LionBaggage struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}

type LionStop struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}

func (f *LionFlight) Validate() error {
	locDep, err := time.LoadLocation(f.Schedule.DepartureTimezone)
	if err != nil {
		return fmt.Errorf("[%s] invalid departure timezone: %s", f.ID, f.Schedule.DepartureTimezone)
	}

	locArr, err := time.LoadLocation(f.Schedule.ArrivalTimezone)
	if err != nil {
		return fmt.Errorf("[%s] invalid arrival timezone: %s", f.ID, f.Schedule.ArrivalTimezone)
	}

	const layout = "2006-01-02T15:04:05"
	_, err = time.ParseInLocation(layout, f.Schedule.Departure, locDep)
	if err != nil {
		return fmt.Errorf("[%s] failed to parse departure time: %w", f.ID, err)
	}

	_, err = time.ParseInLocation(layout, f.Schedule.Arrival, locArr)
	if err != nil {
		return fmt.Errorf("[%s] failed to parse arrival time: %w", f.ID, err)
	}

	if f.Route.From.Code == f.Route.To.Code {
		return fmt.Errorf("[%s] circular route detected: %s to %s",
			f.ID, f.Route.From.Code, f.Route.To.Code)
	}

	if f.Pricing.Total <= 0 {
		return fmt.Errorf("[%s] invalid price: %.2f", f.ID, f.Pricing.Total)
	}

	if f.SeatsLeft < 0 {
		return fmt.Errorf("[%s] negative seat inventory", f.ID)
	}

	return nil
}
