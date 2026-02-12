package entity

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
