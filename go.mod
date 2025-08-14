module github.com/rcglezreyes/go_weather

go 1.22.0

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/labstack/echo/v4 v4.11.4
	github.com/prometheus/client_golang v1.17.0
	github.com/swaggo/echo-swagger v1.4.0
	github.com/swaggo/files v1.0.1
	github.com/swaggo/swag v1.16.2
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240528184218-531527333157
