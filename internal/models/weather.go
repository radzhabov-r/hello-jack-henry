package models

type WeatherRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type WeatherResponse struct {
	Forecast    string `json:"forecast"`
	Temperature string `json:"temperature"`
}

type NWSPointResponse struct {
	Properties struct {
		Forecast         string `json:"forecast"`
		GridId           string `json:"gridId"`
		GridX            int    `json:"gridX"`
		GridY            int    `json:"gridY"`
		Cwa              string `json:"cwa"`
		ForecastOffice   string `json:"forecastOffice"`
		TimeZone         string `json:"timeZone"`
		RelativeLocation struct {
			Properties struct {
				City  string `json:"city"`
				State string `json:"state"`
			} `json:"properties"`
		} `json:"relativeLocation"`
	} `json:"properties"`
}

type NWSForecastResponse struct {
	Properties struct {
		Units             string   `json:"units"`
		ForecastGenerator string   `json:"forecastGenerator"`
		GeneratedAt       string   `json:"generatedAt"`
		UpdateTime        string   `json:"updateTime"`
		ValidTimes        string   `json:"validTimes"`
		Periods           []Period `json:"periods"`
	} `json:"properties"`
}

type Period struct {
	Number                     int    `json:"number"`
	Name                       string `json:"name"`
	StartTime                  string `json:"startTime"`
	EndTime                    string `json:"endTime"`
	IsDaytime                  bool   `json:"isDaytime"`
	Temperature                int    `json:"temperature"`
	TemperatureUnit            string `json:"temperatureUnit"`
	TemperatureTrend           string `json:"temperatureTrend"`
	ProbabilityOfPrecipitation struct {
		UnitCode string `json:"unitCode"`
		Value    int    `json:"value"`
	} `json:"probabilityOfPrecipitation"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	Icon             string `json:"icon"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}
