package router

import (
	"weather-service/internal/handler"

	"github.com/gorilla/mux"
)

func NewRouter(weatherHandler *handler.WeatherHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/weather", weatherHandler.GetWeather).Methods("GET")
	router.HandleFunc("/health", weatherHandler.HealthCheck).Methods("GET") // added just in case, usually we used this for container healthcheck

	return router
}
