package service

import (
	"flight-aggregator/internal/entity"
	"flight-aggregator/internal/service/batikair"
	"flight-aggregator/internal/service/garuda"
	"flight-aggregator/internal/service/lionair"
)

type flightService struct {
	garudaService   garuda.GarudaService
	batikAirService batikair.BatikAirService
	lionAirService  lionair.LionAirService
}

type FlightService interface {
	SearchFlight() ([]entity.LionFlight, error)
}

func NewFlightService(
	garudaService garuda.GarudaService,
	batikAirService batikair.BatikAirService,
	lionAirService lionair.LionAirService,
) FlightService {
	return &flightService{
		garudaService:   garudaService,
		batikAirService: batikAirService,
		lionAirService:  lionAirService,
	}
}

func (f *flightService) SearchFlight() ([]entity.LionFlight, error) {
	// get flight detail
	data, err := f.flightData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *flightService) flightData() ([]entity.LionFlight, error) {
	batik, err := f.lionAirService.GetFlight()
	if err != nil {
		return nil, err
	}

	return batik, nil
}
