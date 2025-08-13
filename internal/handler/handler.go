package handler

import (
	"weather-service/internal/config"
	"weather-service/internal/models"
)

// IWeatherService defines the interface for the weather service
// to allow for easier testing and decoupling.
// I prefer to do interfaces in namespace where they are used, not where they are implemented.
type IWeatherService interface {
	// GetForecastData fetches the weather forecast data for given latitude and longitude.
	GetForecastData(lat float64, lon float64) (*models.NWSForecastResponse, error)
}

func NewWeatherHandler(weatherService IWeatherService, cfg *config.Config) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
		config:         cfg,
	}
}
