package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"weather-service/internal/config"
	"weather-service/internal/models"
)

type WeatherHandler struct {
	weatherService IWeatherService
	config         *config.Config
}

func (wh *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	lat, lon, err := wh.parseAndValidateCoordinates(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	forecastData, err := wh.weatherService.GetForecastData(lat, lon)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	if len(forecastData.Properties.Periods) == 0 {
		http.Error(w, "No forecast periods available", http.StatusInternalServerError)
		return
	}

	todayPeriod, err := wh.findTodaysPeriod(forecastData.Properties.Periods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tempCategory := wh.categorizeTemperature(todayPeriod.Temperature)

	response := &models.WeatherResponse{
		Forecast:    todayPeriod.ShortForecast,
		Temperature: tempCategory,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (wh *WeatherHandler) parseAndValidateCoordinates(r *http.Request) (float64, float64, error) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		return 0, 0, fmt.Errorf("missing required parameters: lat and lon")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude parameter")
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude parameter")
	}

	// out of bounds check
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return 0, 0, fmt.Errorf("invalid coordinates")
	}

	return lat, lon, nil
}

// findTodaysPeriod searches for the most appropriate weather period for today.
// this was very tricky to get right ways of defining "today".
// The search order is as follows:
// - Check for a period named "This Afternoon".
// - Check for any daytime period (IsDaytime == true) that starts today.
// - Check for any period that starts today.
// - Check for any period where the current time falls within its start and end time.
// - Check for the first daytime period.
// - Fallback to the first period if available.
// If none found, return an error indicating no forecast periods are available.
func (wh *WeatherHandler) findTodaysPeriod(periods []models.Period) (models.Period, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// First, look for "This Afternoon" as it's the most common "today" period
	for _, period := range periods {
		if period.Name == "This Afternoon" {
			return period, nil
		}
	}

	// Second, look for any daytime period that starts today
	for _, period := range periods {
		if period.IsDaytime {
			startTime, err := time.Parse(time.RFC3339, period.StartTime)
			if err == nil && startTime.Format("2006-01-02") == today {
				return period, nil
			}
		}
	}

	// Third, look for any period that starts today
	for _, period := range periods {
		startTime, err := time.Parse(time.RFC3339, period.StartTime)
		if err == nil && startTime.Format("2006-01-02") == today {
			return period, nil
		}
	}

	// Fourth, look for any period where current time falls within start/end time
	for _, period := range periods {
		startTime, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			continue
		}
		endTime, err := time.Parse(time.RFC3339, period.EndTime)
		if err != nil {
			continue
		}

		if now.After(startTime) && now.Before(endTime) {
			return period, nil
		}
	}

	// Fallback: look for the first daytime period
	for _, period := range periods {
		if period.IsDaytime {
			return period, nil
		}
	}

	// Last resort: return the first period if available
	if len(periods) > 0 {
		return periods[0], nil
	}

	return models.Period{}, fmt.Errorf("no forecast periods available")
}

func (wh *WeatherHandler) categorizeTemperature(temp int) string {
	if temp >= wh.config.Weather.TemperatureRanges.HotThreshold {
		return "hot"
	} else if temp <= wh.config.Weather.TemperatureRanges.ColdThreshold {
		return "cold"
	}
	return "moderate"
}
