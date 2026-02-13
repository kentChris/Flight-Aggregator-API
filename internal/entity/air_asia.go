package entity

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
