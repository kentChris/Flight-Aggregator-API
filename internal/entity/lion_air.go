package entity

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
