package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"weather-service/internal/models"
)

type WeatherService struct {
	client *http.Client
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		client: &http.Client{
			//we can add whatever setting we want to here to do requests
			Timeout: 10 * time.Second,
		},
	}
}

func (ws *WeatherService) GetForecastData(lat float64, lon float64) (*models.NWSForecastResponse, error) {
	forecastURL, err := ws.getForecastURL(lat, lon)
	if err != nil {
		log.Printf("[ERROR] Failed to get forecast URL for lat=%.4f, lon=%.4f: %v", lat, lon, err)
		return nil, fmt.Errorf("failed to get forecast URL: %w", err)
	}

	return ws.fetchForecast(forecastURL)
}

// getForecastURL retrieves the forecast URL from the NWS API based on latitude and longitude.
// It returns the forecast URL or an error if the request fails.
func (ws *WeatherService) getForecastURL(lat float64, lon float64) (string, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", lat, lon)
	log.Printf("[DEBUG] Calling NWS points API: %s", url)

	resp, err := ws.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NWS API returned status %d", resp.StatusCode)
	}

	var pointResp models.NWSPointResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointResp); err != nil {
		return "", err
	}

	log.Printf("[DEBUG] Location info - %s, %s (Grid: %s %d,%d, Timezone: %s)",
		pointResp.Properties.RelativeLocation.Properties.City,
		pointResp.Properties.RelativeLocation.Properties.State,
		pointResp.Properties.GridId,
		pointResp.Properties.GridX,
		pointResp.Properties.GridY,
		pointResp.Properties.TimeZone)

	return pointResp.Properties.Forecast, nil
}

// fetchForecast retrieves the weather forecast from the NWS API using the provided forecast URL.
// It returns the forecast data or an error if the request fails.
func (ws *WeatherService) fetchForecast(forecastURL string) (*models.NWSForecastResponse, error) {
	log.Printf("[DEBUG] Calling NWS forecast API: %s", forecastURL)

	resp, err := ws.client.Get(forecastURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	var forecastResp models.NWSForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Successfully fetched forecast with %d periods", len(forecastResp.Properties.Periods))

	// Log compact summary of all periods
	//can be commented/removed later, used it for debugging
	//log.Printf("[DEBUG] All periods:")
	//for _, period := range forecastResp.Properties.Periods {
	//	daynight := "Day"
	//	if !period.IsDaytime {
	//		daynight = "Night"
	//	}
	//	log.Printf("[DEBUG]   %d. %s (%s): %d°%s, %s, %d%% precip",
	//		period.Number,
	//		period.Name,
	//		daynight,
	//		period.Temperature,
	//		period.TemperatureUnit,
	//		period.ShortForecast,
	//		period.ProbabilityOfPrecipitation.Value)
	//}

	return &forecastResp, nil
}
