package usecase

import (
	"context"
	"testing"

	"github.com/rcglezreyes/go_weather/internal/pkg/cache"
)

type fakeNWS struct {
	short string
	temp  float64
	err   error
}

func (f fakeNWS) GetToday(ctx context.Context, lat, lon float64) (string, float64, error) {
	return f.short, f.temp, f.err
}

func TestGetTodayForecast_UsesCache(t *testing.T) {
	c := cache.NewTTLCache(cache.Config{TTL: 60, SweepInterval: 10, MaxEntries: 100})
	svc := NewWeatherService(fakeNWS{short: "Sunny", temp: 90}, c)
	res1, err := svc.GetTodayForecast(context.Background(), 1.23456, 7.89012)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := svc.GetTodayForecast(context.Background(), 1.23456, 7.89012)
	if err != nil {
		t.Fatal(err)
	}
	if res1.ShortForecast != res2.ShortForecast {
		t.Fatal("expected cached same result")
	}
	if res2.Category != "hot" {
		t.Fatalf("want hot, got %s", res2.Category)
	}
}

func TestCategorize(t *testing.T) {
	if categorize(85) != "hot" || categorize(60) != "moderate" || categorize(59.9) != "cold" {
		t.Fatal("categorize thresholds failed")
	}
}
