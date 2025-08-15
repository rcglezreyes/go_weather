package grpcadapter

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	weatherv1 "github.com/rcglezreyes/go_weather/api/proto"
	"github.com/rcglezreyes/go_weather/internal/core/ports"
)

type server struct {
	weatherv1.UnimplementedWeatherServiceServer
	svc ports.WeatherService
}

func New(svc ports.WeatherService) *server { return &server{svc: svc} }

func (s *server) GetTodayForecast(ctx context.Context, req *weatherv1.LatLonRequest) (*weatherv1.ForecastReply, error) {
	res, err := s.svc.GetTodayForecast(ctx, req.GetLat(), req.GetLon())
	if err != nil {
		return nil, err
	}
	return &weatherv1.ForecastReply{
		ShortForecast: res.ShortForecast,
		TemperatureF:  res.TemperatureF,
		Category:      res.Category,
	}, nil
}

func tlsConfigFromEnv() (grpc.ServerOption, bool, error) {
	certFile := os.Getenv("GRPC_TLS_CERT")
	keyFile := os.Getenv("GRPC_TLS_KEY")
	if certFile == "" || keyFile == "" {
		return nil, false, nil
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, false, err
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	if ca := os.Getenv("GRPC_TLS_CA"); ca != "" {
		b, err := os.ReadFile(ca)
		if err != nil {
			return nil, false, err
		}
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(b); !ok {
			return nil, false, errors.New("invalid CA")
		}
		tlsCfg.ClientCAs = pool
	}
	return grpc.Creds(credentials.NewTLS(tlsCfg)), true, nil
}

func Run(addr string, s ports.WeatherService) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	//Interceptors (recovery, prometheus)
	uInts := []grpc.UnaryServerInterceptor{
		recovery.UnaryServerInterceptor(),
		grpc_prom.UnaryServerInterceptor,
	}

	tlsOpt, hasTLS, err := tlsConfigFromEnv()
	if err != nil {
		return nil, err
	}
	var opts []grpc.ServerOption
	if hasTLS {
		opts = append(opts, tlsOpt)
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(uInts...))

	gs := grpc.NewServer(opts...)

	weatherv1.RegisterWeatherServiceServer(gs, New(s))

	// Health + Reflection
	hs := health.NewServer()
	healthpb.RegisterHealthServer(gs, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(gs)

	// Prometheus metrics
	grpc_prom.Register(gs)
	grpc_prom.EnableHandlingTimeHistogram()

	go func() { _ = gs.Serve(lis) }()
	return gs, nil
}
