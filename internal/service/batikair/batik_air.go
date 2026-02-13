package batikair

import (
	"encoding/json"
	"flight-aggregator/internal/entity"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type BatikAirResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Results []entity.BatikFlight `json:"results"`
}

type batikAirService struct {
	filePath string
}

type BatikAirService interface {
	GetFlight() ([]entity.BatikFlight, error)
}

func NewBatikAirService(path string) BatikAirService {
	return &batikAirService{
		filePath: path,
	}
}

// assuming it had 15 December as mock param
func (g *batikAirService) GetFlight() ([]entity.BatikFlight, error) {
	//mock sleep 200 - 400 ms
	delay := time.Duration(rand.Intn(201)+200) * time.Millisecond
	time.Sleep(delay)

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return nil, fmt.Errorf("Batik.getFlight: Failed to get batik data")
	}

	var batikAirResponse BatikAirResponse
	if err := json.Unmarshal(data, &batikAirResponse); err != nil {
		return nil, err
	}

	return batikAirResponse.Results, nil
}
