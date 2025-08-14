package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	initOnce sync.Once

	NWSRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "go_weather", Subsystem: "nws", Name: "requests_total", Help: "Total de requests al cliente NWS",
	})
	NWSRequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "go_weather", Subsystem: "nws", Name: "request_duration_seconds", Help: "Duración de requests NWS",
		Buckets: prometheus.DefBuckets,
	})
)

// register safely handles AlreadyRegisteredError so Init() puede llamarse más de una vez sin pánico.
func register(c prometheus.Collector) {
	if err := prometheus.Register(c); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			// ignoramos: ya estaba registrado
			return
		}
		// para otros errores (muy raro), hacemos panic como MustRegister haría
		panic(err)
	}
}

// Init registra collectors por defecto + métricas custom de forma idempotente
func Init() {
	initOnce.Do(func() {
		register(collectors.NewGoCollector())
		register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		register(NWSRequestsTotal)
		register(NWSRequestDuration)
	})
}
