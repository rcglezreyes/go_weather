MODULE=github.com/rcglezreyes/go_weather
PROTO=api/proto/weather.proto

.PHONY: proto swag run build docker test certs compose-up compose-down

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO)

swag:
	swag init -g cmd/server/main.go

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

docker:
	docker build -t go_weather:latest .

docker-prune:
	docker system prune -f

test:
	go test ./...

certs:
	openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -subj "/CN=localhost" -days 365

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down -v

compose-build-app:
	docker compose build app

compose-up-app:
	docker compose up -d app

compose-down-app:
	docker compose down -v app

compose-build-prometheus:
	docker compose build prometheus

compose-up-prometheus:
	docker compose up -d prometheus

compose-down-prometheus:
	docker compose down -v prometheus