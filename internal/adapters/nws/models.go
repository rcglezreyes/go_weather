package nws

// points: top-level fields
type pointsResp struct {
	Forecast       string `json:"forecast"`
	ForecastHourly string `json:"forecastHourly"`
	GridID         string `json:"gridId"`
	GridX          int    `json:"gridX"`
	GridY          int    `json:"gridY"`
}

// periods: forecast period info
type forecastPeriod struct {
	Name            string  `json:"name"`
	Temperature     float64 `json:"temperature"`
	TemperatureUnit string  `json:"temperatureUnit"`
	ShortForecast   string  `json:"shortForecast"`
}

// forecastTop: top-level fields for forecast response
type forecastTop struct {
	Periods []forecastPeriod `json:"periods"`
}

// forecastWithProps: full forecast response with properties
type forecastWithProps struct {
	Properties forecastTop `json:"properties"`
}
