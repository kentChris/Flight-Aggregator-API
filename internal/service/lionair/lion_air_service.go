package lionair

import (
	"encoding/json"
	"flight-aggregator/internal/entity"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type LionResponse struct {
	Success bool     `json:"success"`
	Data    LionData `json:"data"`
}

type LionData struct {
	AvailableFlights []entity.LionFlight `json:"available_flights"`
}

type lionAirService struct {
	filePath string
}

type LionAirService interface {
	GetFlight() ([]entity.LionFlight, error)
}

func NewLionAirService(path string) LionAirService {
	return &lionAirService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (g *lionAirService) GetFlight() ([]entity.LionFlight, error) {
	// ToDo: add utils function
	//mock sleep 100 - 200 ms
	delay := time.Duration(rand.Intn(101)+100) * time.Millisecond
	time.Sleep(delay)

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return nil, fmt.Errorf("Garuda.getFlight: Failed to get garuda data")
	}

	var response LionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Data.AvailableFlights, nil
}
