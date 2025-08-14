package nws

type pointsResp struct {
	Properties struct {
		Forecast       string `json:"forecast"`
		ForecastHourly string `json:"forecastHourly"`
	} `json:"properties"`
}

type forecastResp struct {
	Properties struct {
		Periods []struct {
			Name            string  `json:"name"`
			Temperature     float64 `json:"temperature"`
			TemperatureUnit string  `json:"temperatureUnit"`
			ShortForecast   string  `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}
