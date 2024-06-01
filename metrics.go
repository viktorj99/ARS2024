package main

import (
	"fmt"
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
		[]string{"status"})

	httpMethodCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_method",
			Help: "HTTP method used.",
		},
		[]string{"method"})

	httpResponseTimeSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time",
			Help: "Duration of the HTTP request.",
		},
		[]string{"endpoint"})

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)

	metricsList = []prometheus.Collector{
		httpHits,
		httpStatusCounter,
		httpMethodCounter,
		httpResponseTimeSeconds,
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

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := &responseWriter{w, http.StatusOK}

		timer := prometheus.NewTimer(httpResponseTimeSeconds.WithLabelValues(fmt.Sprintf("%s %s", r.Method, path))) // creates a time series histogram
		f(rw, r)                                                                                                    // original function call
		timer.ObserveDuration()                                                                                     // needs to wrap around the call so that the data collected is as accurate as possible

		httpHits.Inc()                                                                        // basic HTTP hits counter
		httpStatusCounter.WithLabelValues(strconv.Itoa(rw.statusCode)).Inc()                  // labels values based on status code
		httpMethodCounter.WithLabelValues(r.Method).Inc()                                     // labels values based on HTTP method used
		httpRequestsTotal.With(prometheus.Labels{"method": r.Method, "endpoint": path}).Inc() // increment the http_requests_total counter
	}
}
