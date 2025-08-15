package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	initOnce sync.Once

	NWSRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "go_weather", Subsystem: "nws", Name: "requests_total", Help: "Total requests to client NWS",
	})
	NWSRequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "go_weather", Subsystem: "nws", Name: "request_duration_seconds", Help: "Duration of requests to client NWS",
		Buckets: prometheus.DefBuckets,
	})
)

func register(c prometheus.Collector) {
	if err := prometheus.Register(c); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return
		}
		panic(err)
	}
}

func Init() {
	initOnce.Do(func() {
		register(collectors.NewGoCollector())
		register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		register(NWSRequestsTotal)
		register(NWSRequestDuration)
	})
}
