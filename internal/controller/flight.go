package controller

import (
	"context"
	"encoding/json"
	logger "flight-aggregator/internal/common"
	"flight-aggregator/internal/entity"
	"flight-aggregator/internal/service"
)

type FlightController struct {
	flightSerivice service.FlightService
	logger         *logger.Logger
}

func NewFlightController(flightService service.FlightService) FlightController {
	// this should change to gin http handler instead of returning for mock
	return FlightController{
		flightSerivice: flightService,
		logger:         logger.Init(),
	}
}

func (f *FlightController) SearchFlightData() {
	f.logger.Info("Initialize SearchFlightData")

	context := context.Background()

	// stop := 0

	// MOCK request
	req := entity.SearchRequest{
		Origin:        "CGK",
		Destination:   "SBY",
		DepartureDate: "2025-12-15",
		Passanger:     1,
		CabinClass:    "economy",
		// 		Airlines
		Airlines: []string{entity.GARUDA, entity.LIONAIR},
		// 		Price range
		// PriceMin: 400000,
		// PriceMax: 600000,
		// 		number of stops
		// MaxStops: &stop,
		// 		Duration
		// MaxDuration: 110,
		//		TIME
		// MinDepTime: "03:00",
		// MaxDepTime: "12:00",

		// 			Sort By
		// SortOrder: "ASC",
		// SortBy:    "Price",
	}

	result, err := f.flightSerivice.SearchFlight(context, req)
	if err != nil {
		f.logger.Error(err)
	}

	byteData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		f.logger.Error("Failed to marshal:", err)
		return
	}
	f.logger.Info("result: ", string(byteData))
}
