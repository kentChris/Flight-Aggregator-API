package service

import (
	"flight-aggregator/internal/entity"
	"flight-aggregator/internal/service/airasia"
	"flight-aggregator/internal/service/batikair"
	"flight-aggregator/internal/service/garuda"
	"flight-aggregator/internal/service/lionair"
)

type flightService struct {
	garudaService   garuda.GarudaService
	batikAirService batikair.BatikAirService
	lionAirService  lionair.LionAirService
	airAsiaService  airasia.AirAsiaService
}

type FlightService interface {
	SearchFlight() ([]entity.Flight, error)
}

func NewFlightService(
	garudaService garuda.GarudaService,
	batikAirService batikair.BatikAirService,
	lionAirService lionair.LionAirService,
	airAsiaService airasia.AirAsiaService,
) FlightService {
	return &flightService{
		garudaService:   garudaService,
		batikAirService: batikAirService,
		lionAirService:  lionAirService,
		airAsiaService:  airAsiaService,
	}
}

func (f *flightService) SearchFlight() ([]entity.Flight, error) {
	// get flight detail
	data, err := f.flightData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *flightService) flightData() ([]entity.Flight, error) {
	batik, err := f.garudaService.GetFlight()
	if err != nil {
		return nil, err
	}

	return batik, nil
}
