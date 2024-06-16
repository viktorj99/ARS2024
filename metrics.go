package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

var (
	mutex = &sync.Mutex{}

	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ars_requests_total",
			Help: "Total number of requests for the past 24 hours.",
		},
		[]string{"method", "endpoint"},
	)
	successfulRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ars_successful_requests_total",
			Help: "Total number of successful http hits for the past 24 hours.",
		},
		[]string{"method", "endpoint"},
	)
	failedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ars_failed_requests_total",
			Help: "Total number of failed http hits for the past 24 hours.",
		},
		[]string{"method", "endpoint"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ars_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) of http requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	requestRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ars_request_rate",
			Help: "Request rate per endpoint for the past 24 hours.",
		},
		[]string{"method", "endpoint"},
	)

	// Map to store request timestamps for request rate calculation
	requestTimestamps = make(map[string]time.Time)
)

func init() {
	// Register metrics that will be exposed.
	prometheus.MustRegister(totalRequests, successfulRequests, failedRequests, requestDuration, requestRate)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// Count is a middleware function to count HTTP requests, measure request duration,
// and calculate request rate.
func Count(f http.HandlerFunc, method, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment total request count
		mutex.Lock()
		totalRequests.WithLabelValues(method, endpoint).Inc()
		mutex.Unlock()

		// Execute the handler
		rr := &responseRecorder{w, http.StatusOK}
		f(rr, r)

		// Measure request duration
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(method, endpoint).Observe(duration)

		// Increment successful or failed request count based on status code
		if statusCode := rr.statusCode; statusCode >= 200 && statusCode < 400 {
			mutex.Lock()
			successfulRequests.WithLabelValues(method, endpoint).Inc()
			mutex.Unlock()
		} else {
			mutex.Lock()
			failedRequests.WithLabelValues(method, endpoint).Inc()
			mutex.Unlock()
		}

		// Update request rate
		updateRequestRate(method, endpoint)
	}
}

// updateRequestRate updates the request rate metric for a given method and endpoint.
func updateRequestRate(method, endpoint string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Store current timestamp for the endpoint
	key := method + "_" + endpoint
	requestTimestamps[key] = time.Now()

	// Calculate request rate for the past 24 hours
	var count float64
	currentTime := time.Now()
	for _, timestamp := range requestTimestamps {
		if currentTime.Sub(timestamp) <= 24*time.Hour {
			count++
		}
	}
	requestRate.WithLabelValues(method, endpoint).Set(count)
}
