package MyMetrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_http_requests_total",
			Help: "Total HTTP requests",
		}, []string{"method", "path", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_http_request_duration_seconds",
			Help:    "HTTP request durations in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"},
	)
	InFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_http_in_flight_requests",
			Help: "In-flight HTTP requests",
		},
	)
)

func RegisterCollectors() {
	prometheus.MustRegister(
		RequestCount,
		RequestDuration,
		InFlight,
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)
}

func InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InFlight.Inc()
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rr, r)
		duration := time.Since(start).Seconds()

		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		RequestCount.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rr.status)).Inc()
		InFlight.Dec()
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func PromHandler() http.Handler {
	return promhttp.Handler()
}
