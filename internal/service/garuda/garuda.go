package garuda

import (
	"encoding/json"
	"flight-aggregator/internal/entity"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type GarudaResponse struct {
	Status  string                `json:"status"`
	Flights []entity.GarudaFlight `json:"flights"`
}

type garudaService struct {
	filePath string
}

type GarudaService interface {
	GetFlight() ([]entity.GarudaFlight, error)
}

func NewGarudaService(path string) GarudaService {
	return &garudaService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (g *garudaService) GetFlight() ([]entity.GarudaFlight, error) {
	//mock sleep 50 - 100 ms
	delay := time.Duration(rand.Intn(51)+50) * time.Second
	time.Sleep(delay)

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return nil, fmt.Errorf("Garuda.getFlight: Failed to get garuda data")
	}

	var garudaResponse GarudaResponse
	if err := json.Unmarshal(data, &garudaResponse); err != nil {
		return nil, err
	}

	return garudaResponse.Flights, nil
}
