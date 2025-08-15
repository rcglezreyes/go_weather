# go_weather — Echo + Swagger, gRPC, TLS, Docker, Concurrency, Prometheus

## Execute
With ```make``` (Makefile):
```bash
make install-tools

make check-tools

make download

make proto

make swag

make tidy

make vendor

make test

make coverage

make build

make run

# App REST:     http://localhost:8080/swagger/index.html

```

Without ```make```:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest 

go install github.com/swaggo/swag/cmd/swag@latest

go mod download

protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO)

swag init -g cmd/server/main.go

go mod tidy

go mod vendor

go test ./...

go test -coverprofile=coverage.out ./...

go tool cover -html=coverage.out -o coverage.html

go build -o server ./cmd/server

go run ./cmd/server

# App REST:     http://localhost:8080/swagger/index.html

```

## Execute (Docker & Prometheus)

With ```make``` (Makefile):
```bash
make docker
make compose-up

# Prometheus UI: http://localhost:9092  (scrapea /metrics de la app)
# App REST:     http://localhost:8081/swagger/index.html

```
Without ```make```:
```bash
docker build -t go_weather:latest .
docker compose up -d --build

# Prometheus UI: http://localhost:9092  (scrapea /metrics de la app)
# App REST:     http://localhost:8081/swagger/index.html

```

## gRPC test
```bash
grpcurl -plaintext localhost:9090 list weather.v1.WeatherService
grpcurl -plaintext -d '{"lat":38.8894,"lon":-77.0352}' \
  localhost:9090 weather.v1.WeatherService/GetTodayForecast
```

## Endpoints
- REST: `GET /api/v1/forecast?lat={lat}&lon={lon}`
- Health: `/healthz`, `/readyz`
- Metrics (Prometheus): `/metrics`
- gRPC: `weather.v1.WeatherService/GetTodayForecast`

## Exposed metrics
- **HTTP**: `/metrics` includes `go_*`, `process_*`, and custom metrics:
  - `go_weather_nws_requests_total`
  - `go_weather_nws_request_duration_seconds`
- **gRPC**: `go_grpc_*` (latency/throughput) vía `go-grpc-prometheus`.

## Some data for tests in Swagger
New York, NY, USA:
```json
{
  "latitude": 40.7128,
  "longitude": -74.0060,
}
```
Los Angeles, CA, US:
```json
{
  "latitude": 34.0522,
  "longitude": -118.2437,
}
```

Montgomery, AL, US:
```json
{
  "latitude": 32.3668,
  "longitude": -86.3000,
}
```

Juneau, AK, US:
```json
{
  "latitude": 58.3019,
  "longitude": -134.4197,
}
```

Phoenix, AZ, US:
```json
{
  "latitude": 33.4484,
  "longitude": -112.0740,
}
```

Little Rock, AR, US:
```json
{
  "latitude": 34.7465,
  "longitude": -92.2896,
}
```

Sacramento, CA, US:
```json
{
  "latitude": 38.5816,
  "longitude": -121.4944,
}
```

Denver, CO, US:
```json
{
  "latitude": 39.7392,
  "longitude": -104.9903,
}
```

Hartford, CT, US:
```json
{
  "latitude": 41.7658,
  "longitude": -72.6734,
}
```

Dover, DE, US:
```json
{
  "latitude": 39.1582,
  "longitude": -75.5244,
}
```

Tallahassee, FL, US:
```json
{
  "latitude": 30.4383,
  "longitude": -84.2807,
}
```

Atlanta, GA, US:
```json
{
  "latitude": 33.7490,
  "longitude": -84.3880,
}
```

Honolulu, HI, US:
```json
{
  "latitude": 21.3069,
  "longitude": -157.8583,
}
```

Boise, ID, US:
```json
{
  "latitude": 43.6150,
  "longitude": -116.2023,
}
```

Springfield, IL, US:
```json
{
  "latitude": 39.7817,
  "longitude": -89.6501,
}
```

Indianapolis, IN, US:
```json
{
  "latitude": 39.7684,
  "longitude": -86.1581,
}
```

Des Moines, IA, US:
```json
{
  "latitude": 41.5868,
  "longitude": -93.6250,
}
```

Topeka, KS, US:
```json
{
  "latitude": 39.0473,
  "longitude": -95.6752,
}
```

Frankfort, KY, US:
```json
{
  "latitude": 38.2009,
  "longitude": -84.8733,
}
```

Baton Rouge, LA, US:
```json
{
  "latitude": 30.4515,
  "longitude": -91.1871,
}
```

Augusta, ME, US:
```json
{
  "latitude": 44.3106,
  "longitude": -69.7795,
}
```

Annapolis, MD, US:
```json
{
  "latitude": 38.9784,
  "longitude": -76.4922,
}
```

Boston, MA, US:
```json
{
  "latitude": 42.3601,
  "longitude": -71.0589,
}
```

Lansing, MI, US:
```json
{
  "latitude": 42.7325,
  "longitude": -84.5555,
}
```

Saint Paul, MN, US:
```json
{
  "latitude": 44.9537,
  "longitude": -93.0900,
}
```

Jackson, MS, US:
```json
{
  "latitude": 32.2988,
  "longitude": -90.1848,
}
```

Jefferson City, MO, US:
```json
{
  "latitude": 38.5767,
  "longitude": -92.1735,
}
```

Helena, MT, US:
```json
{
  "latitude": 46.5884,
  "longitude": -112.0245,
}
```

Lincoln, NE, US:
```json
{
  "latitude": 40.8136,
  "longitude": -96.7026,
}
```

Carson City, NV, US:
```json
{
  "latitude": 39.1638,
  "longitude": -119.7674,
}
```

Concord, NH, US:
```json
{
  "latitude": 43.2081,
  "longitude": -71.5376,
}
```

Trenton, NJ, US:
```json
{
  "latitude": 40.2206,
  "longitude": -74.7597,
}
```

Santa Fe, NM, US:
```json
{
  "latitude": 35.6870,
  "longitude": -105.9378,
}
```

Albany, NY, US:
```json
{
  "latitude": 42.6526,
  "longitude": -73.7562,
}
```

Raleigh, NC, US:
```json
{
  "latitude": 35.7796,
  "longitude": -78.6382,
}
```

Bismarck, ND, US:
```json
{
  "latitude": 46.8083,
  "longitude": -100.7837,
}
```

Columbus, OH, US:
```json
{
  "latitude": 39.9612,
  "longitude": -82.9988,
}
```

Oklahoma City, OK, US:
```json
{
  "latitude": 35.4676,
  "longitude": -97.5164,
}
```

Salem, OR, US:
```json
{
  "latitude": 44.9429,
  "longitude": -123.0351,
}
```

Harrisburg, PA, US:
```json
{
  "latitude": 40.2732,
  "longitude": -76.8867,
}
```

Providence, RI, US:
```json
{
  "latitude": 41.8240,
  "longitude": -71.4128,
}
```

Columbia, SC, US:
```json
{
  "latitude": 34.0007,
  "longitude": -81.0348,
}
```

Pierre, SD, US:
```json
{
  "latitude": 44.3683,
  "longitude": -100.3510,
}
```

Nashville, TN, US:
```json
{
  "latitude": 36.1627,
  "longitude": -86.7816,
}
```

Austin, TX, US:
```json
{
  "latitude": 30.2672,
  "longitude": -97.7431,
}
```

Salt Lake City, UT, US:
```json
{
  "latitude": 40.7608,
  "longitude": -111.8910,
}
```

Montpelier, VT, US:
```json
{
  "latitude": 44.2601,
  "longitude": -72.5754,
}
```

Richmond, VA, US:
```json
{
  "latitude": 37.5407,
  "longitude": -77.4360,
}
```

Olympia, WA, US:
```json
{
  "latitude": 47.0379,
  "longitude": -122.9007,
}
```

Charleston, WV, US:
```json
{
  "latitude": 38.3498,
  "longitude": -81.6326,
}
```

Madison, WI, US:
```json
{
  "latitude": 43.0731,
  "longitude": -89.4012,
}
```

Cheyenne, WY, US:
```json
{
  "latitude": 41.1400,
  "longitude": -104.8202,
}
```

Washington, DC, US:
```json
{
  "latitude": 38.9072,
  "longitude": -77.0369,
}
```

Mexico City, Mexico:
```json
{
  "latitude": 19.4326,
  "longitude": -99.1332,
}
```

Madrid, Spain:
```json
{
  "latitude": 40.4168,
  "longitude": -3.7038,
}
```

Buenos Aires, Argentina:
```json
{
  "latitude": -34.6037,
  "longitude": -58.3816,
}
```
