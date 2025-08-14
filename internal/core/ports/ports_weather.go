package ports

import (
	"context"

	"github.com/rcglezreyes/go_weather/internal/core/domain"
)

type NWSClient interface {
	GetToday(ctx context.Context, lat, lon float64) (string, float64, error)
}

type WeatherService interface {
	GetTodayForecast(ctx context.Context, lat, lon float64) (domain.TodayForecast, error)
}
