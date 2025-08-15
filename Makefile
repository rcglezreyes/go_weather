MODULE=github.com/rcglezreyes/go_weather
PROTO=api/proto/weather.proto

# ===== Installation for binaries =====
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell go env GOPATH)/bin
endif

SWAG                := $(GOBIN)/swag
PROTOC_GEN_GO       := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC  := $(GOBIN)/protoc-gen-go-grpc
PROTOC              := $(shell command -v protoc 2>/dev/null)

# Add GOBIN to PATH for this Makefile execution
export PATH := $(PATH):$(GOBIN)

ifneq (,$(wildcard /bin/bash))
  SHELL := /bin/bash
else
  SHELL := /bin/sh
endif

.PHONY: proto swag swag-run run build docker test certs \
        compose-up compose-down compose-build-app compose-up-app compose-down-app \
        compose-build-prometheus compose-up-prometheus compose-down-prometheus \
        install-tools check-tools download tidy vendor docker-prune

# ---- Tools installation ----
install-tools: $(SWAG) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	@echo "✅ Tools ready in $(GOBIN)"

$(SWAG):
	@echo "⬇️  Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest

$(PROTOC_GEN_GO):
	@echo "⬇️  Installing protoc-gen-go..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(PROTOC_GEN_GO_GRPC):
	@echo "⬇️  Installing protoc-gen-go-grpc..."
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

check-tools:
	@echo "swag:               $$($(SWAG) --version 2>/dev/null || echo 'missing')"
	@echo "protoc:             $(if $(PROTOC),$(shell protoc --version),missing)"
	@echo "protoc-gen-go:      $$($(PROTOC_GEN_GO) --version 2>/dev/null || echo 'installed (no --version)')"
	@echo "protoc-gen-go-grpc: $$($(PROTOC_GEN_GO_GRPC) --version 2>/dev/null || echo 'installed (no --version)')"

download:
	go mod download

tidy:
	go mod tidy

vendor:
	go mod vendor

# ---- Verify if protoc is installed ----
ensure-protoc:
	@if [ -z "$(PROTOC)" ]; then \
		echo "❌ Missing 'protoc' compiler."; \
		echo "   Install it first (e.g. macOS: 'brew install protobuf', Ubuntu/Debian: 'sudo apt-get install -y protobuf-compiler', Alpine: 'apk add protobuf')."; \
		exit 1; \
	fi

# ---- Protobuf ----
proto: ensure-protoc install-tools
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO)

# ---- Swagger / OpenAPI ----
swag: install-tools
	$(SWAG) init -g cmd/server/main.go -o ./docs

# To regenerate docs (if needed)
swag-run:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o ./docs

# ---- App ----
run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html


# ---- Docker ----
docker:
	docker build -t go_weather:latest .

docker-prune:
	docker system prune -f

# ---- TLS dev cert ----
certs:
	openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -subj "/CN=localhost" -days 365

# ---- Compose ----
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
