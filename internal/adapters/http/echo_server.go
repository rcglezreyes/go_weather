package httpadapter

import (
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/rcglezreyes/go_weather/internal/adapters/http/handlers"
	apikey "github.com/rcglezreyes/go_weather/internal/adapters/http/middleware/apikey"
	"github.com/rcglezreyes/go_weather/internal/core/ports"
)

func NewEchoServer(svc ports.WeatherService) *echo.Echo {
	e := echo.New()

	//Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(50)))
	e.Use(apikey.OptionalCheckerFromEnv())

	// Health
	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(200) })
	e.GET("/readyz", func(c echo.Context) error { return c.NoContent(200) })

	// Prometheus metrics
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// API
	h := handlers.NewWeatherHandler(svc)
	v1 := e.Group("/api/v1")
	v1.GET("/forecast", h.GetTodayForecast)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e
}
