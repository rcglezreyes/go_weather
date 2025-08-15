package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/rcglezreyes/go_weather/docs" // swagger (si generas con swag)

	grpcadapter "github.com/rcglezreyes/go_weather/internal/adapters/grpc"
	httpadapter "github.com/rcglezreyes/go_weather/internal/adapters/http"
	"github.com/rcglezreyes/go_weather/internal/adapters/nws"
	"github.com/rcglezreyes/go_weather/internal/core/usecase"
	"github.com/rcglezreyes/go_weather/internal/pkg/cache"
	"github.com/rcglezreyes/go_weather/observability/metrics"
)

// @title Weather Forecast API
// @version 1.0
// @description Returns today's short forecast and a temperature category (hot/moderate/cold) using NWS.
// @BasePath /api/v1
func main() {
	httpPort := flag.String("http-port", getenvDefault("PORT", "8080"), "HTTP port")
	grpcPort := flag.String("grpc-port", getenvDefault("GRPC_PORT", "9090"), "gRPC port")
	flag.Parse()

	metrics.Init()

	// Cache: TTL & janitor
	c := cache.NewTTLCache(cache.Config{TTL: 300 /*s*/, SweepInterval: 60 /*s*/, MaxEntries: 5000})

	// Adapters + use case
	nwsClient := nws.NewNWSClient()
	svc := usecase.NewWeatherService(nwsClient, c)

	// gRPC (with Prometheus)
	if _, err := grpcadapter.Run(":"+*grpcPort, svc); err != nil {
		log.Fatalf("gRPC: %v", err)
	}
	log.Printf("gRPC listening on :%s", *grpcPort)

	// HTTP (Echo + Swagger + /metrics)
	e := httpadapter.NewEchoServer(svc)
	log.Printf("HTTP listening on :%s", *httpPort)
	if err := e.Start(":" + *httpPort); err != nil {
		log.Fatal(err)
	}
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
