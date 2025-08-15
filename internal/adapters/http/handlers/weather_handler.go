package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/rcglezreyes/go_weather/internal/core/ports"
)

type WeatherHandler struct{ svc ports.WeatherService }

func NewWeatherHandler(svc ports.WeatherService) *WeatherHandler { return &WeatherHandler{svc: svc} }

// GetTodayForecast godoc
// @Summary Get today's short forecast and temperature category
// @Description Returns today's short forecast and temperature category using NWS
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Produce json
// @Success 200 {object} ForecastResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /forecast [get]
func (h *WeatherHandler) GetTodayForecast(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	latStr := c.QueryParam("lat")
	lonStr := c.QueryParam("lon")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid lat"})
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid lon",
		})
	}

	res, err := h.svc.GetTodayForecast(c.Request().Context(), lat, lon)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadGateway, ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ForecastResponse{
		ShortForecast: res.ShortForecast,
		TemperatureF:  res.TemperatureF,
		Category:      res.Category,
	})
}

type ForecastResponse struct {
	ShortForecast string  `json:"shortForecast"`
	TemperatureF  float64 `json:"temperatureF"`
	Category      string  `json:"category"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
