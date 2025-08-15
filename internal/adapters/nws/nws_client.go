package nws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rcglezreyes/go_weather/internal/pkg/httpclient"
	obs "github.com/rcglezreyes/go_weather/observability/metrics"
)

type Client struct{ http *http.Client }

func NewNWSClient() *Client {
	return &Client{
		http: httpclient.New("go_weather/1.0 (contact: rcglezreyes@gmail.com)"),
	}
}

func (c *Client) GetToday(ctx context.Context, lat, lon float64) (string, float64, error) {
	start := time.Now()
	defer func() { obs.NWSRequestsTotal.Inc() }()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon)
	resp, err := c.doNWS(ctx, http.MethodGet, pointsURL)
	if err != nil {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, fmt.Errorf("points status: %d body: %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB
	if err != nil {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, fmt.Errorf("read points body: %w", err)
	}

	var p pointsResp
	if err := json.Unmarshal(body, &p); err != nil {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, fmt.Errorf("unmarshal points: %w; body=%s", err, string(body))
	}

	if p.Forecast == "" && p.ForecastHourly == "" {
		obs.NWSRequestDuration.Observe(time.Since(start).Seconds())
		return "", 0, errors.New("no forecast URLs from points endpoint (after fallback)")
	}

	var rForecast, rHourly []period
	var errForecast, errHourly error

	g, ctx2 := errgroup.WithContext(ctx)

	if p.Forecast != "" {
		url := p.Forecast
		g.Go(func() error {
			pr, err := c.fetchPeriods(ctx2, url)
			if err == nil {
				rForecast = pr
			} else {
				errForecast = err
			}
			return nil
		})
	}
	if p.ForecastHourly != "" {
		url := p.ForecastHourly
		g.Go(func() error {
			pr, err := c.fetchPeriods(ctx2, url)
			if err == nil {
				rHourly = pr
			} else {
				errHourly = err
			}
			return nil
		})
	}
	_ = g.Wait()

	if len(rForecast) == 0 && len(rHourly) == 0 {
		return "", 0, fmt.Errorf("no forecast periods: dailyErr=%v hourlyErr=%v", errForecast, errHourly)
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
	resp, err := c.doNWS(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("forecast status: %d body: %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("read forecast body: %w", err)
	}

	// 1) periods in top-level
	var ft forecastTop
	if err := json.Unmarshal(body, &ft); err == nil && len(ft.Periods) > 0 {
		out := make([]period, 0, len(ft.Periods))
		for _, pr := range ft.Periods {
			out = append(out, period{
				Name:  pr.Name,
				Temp:  pr.Temperature,
				Unit:  pr.TemperatureUnit,
				Short: pr.ShortForecast,
			})
		}
		return out, nil
	}

	// 2) fallback: properties.periods
	var fp forecastWithProps
	if err := json.Unmarshal(body, &fp); err == nil && len(fp.Properties.Periods) > 0 {
		out := make([]period, 0, len(fp.Properties.Periods))
		for _, pr := range fp.Properties.Periods {
			out = append(out, period{
				Name:  pr.Name,
				Temp:  pr.Temperature,
				Unit:  pr.TemperatureUnit,
				Short: pr.ShortForecast,
			})
		}
		return out, nil
	}

	return nil, fmt.Errorf("empty periods from %s; body=%s", url, string(body))
}

func (c *Client) doNWS(ctx context.Context, method, url string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, method, url, nil)
	req.Header.Del("Accept")
	req.Header["Accept"] = []string{"application/geo+json"}
	req.Header.Set("User-Agent", "go_weather/1.0 (contact: rcglezreyes@gmail.com)")
	return c.http.Do(req)
}

func chooseToday(periods []period) (string, float64, bool) {
	for _, pr := range periods {
		n := strings.ToLower(pr.Name)
		if n == "today" || strings.Contains(n, "today") ||
			strings.Contains(n, "this morning") || strings.Contains(n, "this afternoon") ||
			strings.Contains(n, "tonight") {
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
