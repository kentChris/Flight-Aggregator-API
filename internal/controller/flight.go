package controller

import (
	"encoding/json"
	logger "flight-aggregator/internal/common"
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

	result, err := f.flightSerivice.SearchFlight()
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
