package main

import (
	logger "flight-aggregator/internal/common"
	"flight-aggregator/internal/controller"
	"flight-aggregator/internal/redis"
	"flight-aggregator/internal/service"
	"flight-aggregator/internal/service/airasia"
	"flight-aggregator/internal/service/batikair"
	"flight-aggregator/internal/service/garuda"
	"flight-aggregator/internal/service/lionair"
)

func main() {
	log := logger.Init()

	log.Info("Starting the app")

	// Init Service
	garudaService := garuda.NewGarudaService("mock/garuda_indonesia_search_response.json")
	batikAirService := batikair.NewBatikAirService("mock/batik_air_search_response.json")
	lionAirService := lionair.NewLionAirService("mock/lion_air_search_response.json")
	airasia := airasia.NewAirAsiaService("mock/airasia_search_response.json")
	// We can enhance by storing this in env
	// Change this to localhost:6379 for running in local (without docker compose)
	// redisService := redis.NewRedisService("redis:6379", "", 0)
	redisService := redis.NewRedisService("localhost:6379", "", 0)
	flightService := service.NewFlightService(garudaService, batikAirService, lionAirService, airasia, redisService)

	// Init controller
	flightController := controller.NewFlightController(flightService)

	// mock
	mock(flightController)
}

func mock(ctrl controller.FlightController) {
	ctrl.SearchFlightData()
}
