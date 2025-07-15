package handlers

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "platform_http_requests_total",
			Help: "Total HTTP requests by method, path and status",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "platform_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "platform_active_connections",
		Help: "Current number of active HTTP connections",
	})

	appInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "platform_app_info",
			Help: "Application information",
		},
		[]string{"version", "environment"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration, activeConnections, appInfo)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.written {
		r.statusCode = code
		r.written = true
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.written {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		activeConnections.Inc()
		defer activeConnections.Dec()

		timer := prometheus.NewTimer(requestDuration.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		status := strconv.Itoa(recorder.statusCode)
		requestsTotal.WithLabelValues(method, path, status).Inc()
	})
}
