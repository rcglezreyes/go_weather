package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/rcglezreyes/go_weather/internal/core/domain"
	"github.com/rcglezreyes/go_weather/internal/core/ports"
	"github.com/rcglezreyes/go_weather/internal/pkg/cache"
)

type weatherService struct {
	nws   ports.NWSClient
	cache cache.KV
}

func NewWeatherService(nws ports.NWSClient, c cache.KV) ports.WeatherService {
	return &weatherService{nws: nws, cache: c}
}

func (s *weatherService) GetTodayForecast(ctx context.Context, lat, lon float64) (domain.TodayForecast, error) {
	key := cacheKey(lat, lon)
	if v, ok := s.cache.Get(key); ok {
		return v.(domain.TodayForecast), nil
	}

	short, tempF, err := s.nws.GetToday(ctx, lat, lon)
	if err != nil {
		return domain.TodayForecast{}, err
	}

	res := domain.TodayForecast{ShortForecast: short, TemperatureF: tempF, Category: categorize(tempF)}
	s.cache.Set(key, res)
	return res, nil
}

func categorize(tempF float64) string {
	switch {
	case tempF >= 85:
		return "hot"
	case tempF >= 60:
		return "moderate"
	default:
		return "cold"
	}
}

func cacheKey(lat, lon float64) string {
	r := func(x float64) float64 { return math.Round(x*1000) / 1000 }
	return fmt.Sprintf("lat=%.3f:lon=%.3f", r(lat), r(lon))
}
