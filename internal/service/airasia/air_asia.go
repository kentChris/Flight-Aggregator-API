package airasia

import (
	"encoding/json"
	"flight-aggregator/internal/entity"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type AirAsiaResponse struct {
	Status  string                 `json:"status"`
	Flights []entity.AirAsiaFlight `json:"flights"`
}

type airAsiaService struct {
	filePath string
}

type AirAsiaService interface {
	GetFlight() ([]entity.AirAsiaFlight, error)
}

func NewAirAsiaService(path string) AirAsiaService {
	return &airAsiaService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (g *airAsiaService) GetFlight() ([]entity.AirAsiaFlight, error) {
	//mock sleep 50 - 150 ms
	// ToDo: add utils function
	delay := time.Duration(rand.Intn(101)+50) * time.Millisecond
	time.Sleep(delay)

	// ToDo: add this in util
	change := rand.Intn(100)
	if change >= 90 {
		return nil, fmt.Errorf("Mock Fail")
	}

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return nil, fmt.Errorf("airAsia.getFlight: Failed to get Air Asia data")
	}

	var airAsiaResponse AirAsiaResponse
	if err := json.Unmarshal(data, &airAsiaResponse); err != nil {
		return nil, err
	}

	return airAsiaResponse.Flights, nil
}
