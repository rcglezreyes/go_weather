FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git make protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:${PATH}"
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN make proto \
 && swag init -g cmd/server/main.go \
 && go mod tidy \
 && go build -o server ./cmd/server


# --- runtime ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 8080 9090
ENV PORT=8080 GRPC_PORT=9090
CMD ["/app/server"]