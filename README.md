# go_weather — Echo + Swagger, gRPC, TLS, Docker, Concurrency, Prometheus

## Endpoints
- REST: `GET /api/v1/forecast?lat={lat}&lon={lon}`
- Health: `/healthz`, `/readyz`
- Metrics (Prometheus): `/metrics`
- gRPC: `weather.v1.WeatherService/GetTodayForecast`

## Métricas expuestas
- **HTTP**: `/metrics` incluye `go_*`, `process_*`, y métricas custom:
  - `go_weather_nws_requests_total`
  - `go_weather_nws_request_duration_seconds`
- **gRPC**: `go_grpc_*` (latencias/contadores) vía `go-grpc-prometheus`.

## Ejecutar
```bash
go mod tidy
make proto
make swag
make run
```

## Docker & Prometheus
```bash
make docker
make compose-up
# Prometheus UI: http://localhost:9091  (scrapea /metrics de la app)
# App REST:     http://localhost:8080/swagger/index.html
```

## gRPC prueba
```bash
grpcurl -plaintext localhost:9090 list weather.v1.WeatherService
grpcurl -plaintext -d '{"lat":38.8894,"lon":-77.0352}' \
  localhost:9090 weather.v1.WeatherService/GetTodayForecast
```
```
