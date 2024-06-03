package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

var (
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total",
			Help: "Total number of HTTP hits.",
		},
	)

	httpStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of the HTTP response.",
		},
		[]string{"status"},
	)

	httpMethodCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_method",
			Help: "HTTP method used.",
		},
		[]string{"method"},
	)

	httpResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time",
			Help: "Duration of the HTTP request.",
		},
		[]string{"endpoint"},
	)

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status"},
	)

	metricsList = []prometheus.Collector{
		httpHits,
		httpStatusCounter,
		httpMethodCounter,
		httpResponseTime,
		httpRequestsTotal,
	}

	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := &responseWriter{w, http.StatusOK}

		timer := prometheus.NewTimer(httpResponseTime.WithLabelValues(path))
		next.ServeHTTP(rw, r)
		timer.ObserveDuration()

		httpHits.Inc()
		httpStatusCounter.WithLabelValues(strconv.Itoa(rw.statusCode)).Inc()
		httpMethodCounter.WithLabelValues(r.Method).Inc()
		httpRequestsTotal.With(prometheus.Labels{
			"method":   r.Method,
			"endpoint": path,
			"status":   strconv.Itoa(rw.statusCode),
		}).Inc()
	})
}
