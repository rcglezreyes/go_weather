package grpcadapter

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	weatherv1 "github.com/rcglezreyes/go_weather/api/proto"
	"github.com/rcglezreyes/go_weather/internal/core/domain"
)

type fakeSvc struct{}

func (fakeSvc) GetTodayForecast(ctx context.Context, lat, lon float64) (domain.TodayForecast, error) {
	return domain.TodayForecast{ShortForecast: "Sunny", TemperatureF: 75, Category: "moderate"}, nil
}

func dialer(gs *grpc.Server) func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	go func() { _ = gs.Serve(lis) }()
	return func(context.Context, string) (net.Conn, error) { return lis.Dial() }
}

func TestGRPC_GetTodayForecast(t *testing.T) {
	gs := grpc.NewServer()
	weatherv1.RegisterWeatherServiceServer(gs, New(fakeSvc{}))

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(dialer(gs)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	cli := weatherv1.NewWeatherServiceClient(conn)
	got, err := cli.GetTodayForecast(context.Background(), &weatherv1.LatLonRequest{Lat: 1, Lon: 2})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetCategory() != "moderate" {
		t.Fatalf("want moderate, got %s", got.GetCategory())
	}
}
