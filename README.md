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
make docker-prune
make docker
make compose-up

# Prometheus UI: http://localhost:9092  (scrap metrics gRPC)
# App REST:     http://localhost:8081/swagger/index.html

```
Without ```make```:
```bash
docker build -t go_weather:latest .
docker compose up -d --build

# Prometheus UI: http://localhost:9092  (scrap metrics gRPC)
# App REST:     http://localhost:8081/swagger/index.html

```

## Endpoints
- REST: `GET /api/v1/forecast?lat={lat}&lon={lon}`
- Health: `/healthz`, `/readyz`
- Metrics (Prometheus): `/metrics`
- gRPC: `weather.v1.WeatherService/GetTodayForecast` (Must generate certs and declare API KEY as env var)

## Exposed metrics
- **HTTP**: `/metrics` includes `go_*`, `process_*`, and custom metrics:
  - `go_weather_nws_requests_total`
  - `go_weather_nws_request_duration_seconds`
- **gRPC**: `go_grpc_*` (latency/throughput) vía `go-grpc-prometheus`.


## gRPC test (with Docker)
Must generate certs and declare API KEY as env var
```bash
grpcurl -plaintext localhost:9091 list weather.v1.WeatherService
grpcurl -plaintext -d '{"lat":38.8894,"lon":-77.0352}' \
  localhost:9091 weather.v1.WeatherService/GetTodayForecast
```

## gRPC test (without Docker)
Must generate certs and declare API KEY as env var
```bash
grpcurl -plaintext localhost:9090 list weather.v1.WeatherService
grpcurl -plaintext -d '{"lat":38.8894,"lon":-77.0352}' \
  localhost:9090 weather.v1.WeatherService/GetTodayForecast
```


## Some data for tests in Swagger
New York, NY, USA:
```json
{
  "lat": 40.7128,
  "lon": -74.0060,
}
```
Los Angeles, CA, US:
```json
{
  "lat": 34.0522,
  "lon": -118.2437,
}
```

Montgomery, AL, US:
```json
{
  "lat": 32.3668,
  "lon": -86.3000,
}
```

Juneau, AK, US:
```json
{
  "lat": 58.3019,
  "lon": -134.4197,
}
```

Phoenix, AZ, US:
```json
{
  "lat": 33.4484,
  "lon": -112.0740,
}
```

Little Rock, AR, US:
```json
{
  "lat": 34.7465,
  "lon": -92.2896,
}
```

Sacramento, CA, US:
```json
{
  "lat": 38.5816,
  "lon": -121.4944,
}
```

Denver, CO, US:
```json
{
  "lat": 39.7392,
  "lon": -104.9903,
}
```

Hartford, CT, US:
```json
{
  "lat": 41.7658,
  "lon": -72.6734,
}
```

Dover, DE, US:
```json
{
  "lat": 39.1582,
  "lon": -75.5244,
}
```

Tallahassee, FL, US:
```json
{
  "lat": 30.4383,
  "lon": -84.2807,
}
```

Atlanta, GA, US:
```json
{
  "lat": 33.7490,
  "lon": -84.3880,
}
```

Honolulu, HI, US:
```json
{
  "lat": 21.3069,
  "lon": -157.8583,
}
```

Boise, ID, US:
```json
{
  "lat": 43.6150,
  "lon": -116.2023,
}
```

Springfield, IL, US:
```json
{
  "lat": 39.7817,
  "lon": -89.6501,
}
```

Indianapolis, IN, US:
```json
{
  "lat": 39.7684,
  "lon": -86.1581,
}
```

Des Moines, IA, US:
```json
{
  "lat": 41.5868,
  "lon": -93.6250,
}
```

Topeka, KS, US:
```json
{
  "lat": 39.0473,
  "lon": -95.6752,
}
```

Frankfort, KY, US:
```json
{
  "lat": 38.2009,
  "lon": -84.8733,
}
```

Baton Rouge, LA, US:
```json
{
  "lat": 30.4515,
  "lon": -91.1871,
}
```

Augusta, ME, US:
```json
{
  "lat": 44.3106,
  "lon": -69.7795,
}
```

Annapolis, MD, US:
```json
{
  "lat": 38.9784,
  "lon": -76.4922,
}
```

Boston, MA, US:
```json
{
  "lat": 42.3601,
  "lon": -71.0589,
}
```

Lansing, MI, US:
```json
{
  "lat": 42.7325,
  "lon": -84.5555,
}
```

Saint Paul, MN, US:
```json
{
  "lat": 44.9537,
  "lon": -93.0900,
}
```

Jackson, MS, US:
```json
{
  "lat": 32.2988,
  "lon": -90.1848,
}
```

Jefferson City, MO, US:
```json
{
  "lat": 38.5767,
  "lon": -92.1735,
}
```

Helena, MT, US:
```json
{
  "lat": 46.5884,
  "lon": -112.0245,
}
```

Lincoln, NE, US:
```json
{
  "lat": 40.8136,
  "lon": -96.7026,
}
```

Carson City, NV, US:
```json
{
  "lat": 39.1638,
  "lon": -119.7674,
}
```

Concord, NH, US:
```json
{
  "lat": 43.2081,
  "lon": -71.5376,
}
```

Trenton, NJ, US:
```json
{
  "lat": 40.2206,
  "lon": -74.7597,
}
```

Santa Fe, NM, US:
```json
{
  "lat": 35.6870,
  "lon": -105.9378,
}
```

Albany, NY, US:
```json
{
  "lat": 42.6526,
  "lon": -73.7562,
}
```

Raleigh, NC, US:
```json
{
  "lat": 35.7796,
  "lon": -78.6382,
}
```

Bismarck, ND, US:
```json
{
  "lat": 46.8083,
  "lon": -100.7837,
}
```

Columbus, OH, US:
```json
{
  "lat": 39.9612,
  "lon": -82.9988,
}
```

Oklahoma City, OK, US:
```json
{
  "lat": 35.4676,
  "lon": -97.5164,
}
```

Salem, OR, US:
```json
{
  "lat": 44.9429,
  "lon": -123.0351,
}
```

Harrisburg, PA, US:
```json
{
  "lat": 40.2732,
  "lon": -76.8867,
}
```

Providence, RI, US:
```json
{
  "lat": 41.8240,
  "lon": -71.4128,
}
```

Columbia, SC, US:
```json
{
  "lat": 34.0007,
  "lon": -81.0348,
}
```

Pierre, SD, US:
```json
{
  "lat": 44.3683,
  "lon": -100.3510,
}
```

Nashville, TN, US:
```json
{
  "lat": 36.1627,
  "lon": -86.7816,
}
```

Austin, TX, US:
```json
{
  "lat": 30.2672,
  "lon": -97.7431,
}
```

Salt Lake City, UT, US:
```json
{
  "lat": 40.7608,
  "lon": -111.8910,
}
```

Montpelier, VT, US:
```json
{
  "lat": 44.2601,
  "lon": -72.5754,
}
```

Richmond, VA, US:
```json
{
  "lat": 37.5407,
  "lon": -77.4360,
}
```

Olympia, WA, US:
```json
{
  "lat": 47.0379,
  "lon": -122.9007,
}
```

Charleston, WV, US:
```json
{
  "lat": 38.3498,
  "lon": -81.6326,
}
```

Madison, WI, US:
```json
{
  "lat": 43.0731,
  "lon": -89.4012,
}
```

Cheyenne, WY, US:
```json
{
  "lat": 41.1400,
  "lon": -104.8202,
}
```

Washington, DC, US:
```json
{
  "lat": 38.9072,
  "lon": -77.0369,
}
```

Mexico City, Mexico:
```json
{
  "lat": 19.4326,
  "lon": -99.1332,
}
```

Madrid, Spain:
```json
{
  "lat": 40.4168,
  "lon": -3.7038,
}
```

Buenos Aires, Argentina:
```json
{
  "lat": -34.6037,
  "lon": -58.3816,
}
```
