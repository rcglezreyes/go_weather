package nws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rcglezreyes/go_weather/internal/pkg/httpclient"
	obs "github.com/rcglezreyes/go_weather/observability/metrics"
)

type Client struct{ http *http.Client }

func NewNWSClient() *Client {
	return &Client{http: httpclient.New("go_weather/1.0 (contact: rcglezreyes@gmail.com)")}
}

func (c *Client) GetToday(ctx context.Context, lat, lon float64) (string, float64, error) {
	start := time.Now()
	defer func() { obs.NWSRequestsTotal.Inc() }()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, pointsURL, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, fmt.Errorf("points status: %d", resp.StatusCode)
	}
	fmt.Println(resp.Body)
	var p pointsResp
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, err
	}
	if p.Properties.Forecast == "" && p.Properties.ForecastHourly == "" {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, errors.New("no forecast URLs from points endpoint")
	}

	var rForecast, rHourly []period
	g, ctx2 := errgroup.WithContext(ctx)

	if p.Properties.Forecast != "" {
		url := p.Properties.Forecast
		g.Go(func() error {
			pr, err := c.fetchPeriods(ctx2, url)
			if err == nil {
				rForecast = pr
			}
			return err
		})
	}
	if p.Properties.ForecastHourly != "" {
		url := p.Properties.ForecastHourly
		g.Go(func() error {
			pr, err := c.fetchPeriods(ctx2, url)
			if err == nil {
				rHourly = pr
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		if len(rForecast) == 0 && len(rHourly) == 0 {
			obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
			return "", 0, err
		}
	}

	if short, tempF, ok := chooseToday(rForecast); ok {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return short, tempF, nil
	}
	if short, tempF, ok := chooseToday(rHourly); ok {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return short, tempF, nil
	}

	obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
	return "", 0, errors.New("no forecast periods available")
}

type period struct {
	Name  string
	Temp  float64
	Unit  string
	Short string
}

func (c *Client) fetchPeriods(ctx context.Context, url string) ([]period, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("forecast status: %d", resp.StatusCode)
	}

	var f forecastResp
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	out := make([]period, 0, len(f.Properties.Periods))
	for _, pr := range f.Properties.Periods {
		out = append(out, period{Name: pr.Name, Temp: pr.Temperature, Unit: pr.TemperatureUnit, Short: pr.ShortForecast})
	}
	return out, nil
}

func chooseToday(periods []period) (string, float64, bool) {
	for _, pr := range periods {
		n := strings.ToLower(pr.Name)
		if n == "today" || strings.Contains(n, "today") {
			return pr.Short, normalizeF(pr.Temp, pr.Unit), true
		}
	}
	if len(periods) > 0 {
		pr := periods[0]
		return pr.Short, normalizeF(pr.Temp, pr.Unit), true
	}
	return "", 0, false
}

func normalizeF(v float64, unit string) float64 {
	if strings.ToUpper(unit) == "C" {
		return (v * 9.0 / 5.0) + 32.0
	}
	return v
}
