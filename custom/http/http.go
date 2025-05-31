// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package http provides HTTP middleware for metrics collection and monitoring.
// It offers tools to track and measure HTTP request counts, durations, and response codes,
// allowing for detailed monitoring and analysis of HTTP traffic in Go applications.
package http

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type (
	// HTTPMetricsMiddleware defines an interface for HTTP metrics collection middleware.
	// It provides a way to collect and report metrics for HTTP requests, enabling
	// monitoring of API performance, traffic patterns, and error rates.
	HTTPMetricsMiddleware interface {
		// Handler wraps an existing http.Handler with metrics collection.
		// It tracks request counts and durations with attributes for method, URI, and status code,
		// providing detailed insights into HTTP request handling.
		Handler(next http.Handler) http.Handler
	}

	// httpMetricsMiddleware implements the HTTPMetricsMiddleware interface.
	// It uses OpenTelemetry metrics instruments to track HTTP request data.
	httpMetricsMiddleware struct {
		// meter is the OpenTelemetry meter used to create metrics instruments.
		// It serves as the entry point for creating metrics collectors.
		meter metric.Meter

		// requestCounter counts the number of HTTP requests processed.
		// It's used to track traffic volume and patterns over time.
		requestCounter metric.Int64Counter

		// requestDuration measures the duration of HTTP requests.
		// It provides insights into latency and performance characteristics.
		requestDuration metric.Float64Histogram
	}

	// responseWriter wraps an http.ResponseWriter to capture the status code.
	// This allows the middleware to record the final status of the HTTP response
	// for metrics collection.
	responseWriter struct {
		http.ResponseWriter
		statusCode int
	}
)

// NewHTTPMetricsMiddleware creates a new HTTP metrics middleware that collects
// request counts and durations for HTTP requests. It sets up OpenTelemetry
// instruments for tracking request metrics with standardized names and descriptions.
//
// Returns:
//   - An HTTPMetricsMiddleware interface for HTTP metrics collection.
//   - An error if the meter instruments cannot be created.
func NewHTTPMetricsMiddleware() (HTTPMetricsMiddleware, error) {
	// Create a meter with an appropriate instrumentation scope name
	meter := otel.Meter("github.com/goxkit/metrics/custom/http")

	// Create a counter for tracking the total number of HTTP requests
	counter, err := meter.Int64Counter("http.requests", metric.WithDescription("HTTP Requests Counter"))
	if err != nil {
		return nil, err
	}

	// Create a histogram for measuring HTTP request durations
	duration, err := meter.Float64Histogram("http.request.duration", metric.WithDescription("HTTP Request Duration"))
	if err != nil {
		return nil, err
	}

	// Return the configured middleware implementation
	return &httpMetricsMiddleware{
		meter:           meter,
		requestCounter:  counter,
		requestDuration: duration,
	}, nil
}

// Handler wraps an HTTP handler with metrics collection functionality.
// It records the request duration and increments the request counter
// with method, URI, and status code attributes, providing valuable insights
// into API usage patterns, performance characteristics, and error rates.
//
// Parameters:
//   - next: The HTTP handler to wrap with metrics collection.
//
// Returns:
//   - An HTTP handler that collects metrics before calling the wrapped handler.
func (m *httpMetricsMiddleware) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Preserve the request context
		ctx := r.Context()

		// Wrap the response writer to capture the status code
		rw := &responseWriter{w, http.StatusOK}

		// Record the start time for duration calculation
		start := time.Now()

		// Process the request with the wrapped handler
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Record the request duration with method, URI, and status attributes
		m.requestDuration.Record(
			ctx,
			float64(time.Since(start).Nanoseconds()),
			metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.String("uri", r.RequestURI),
				attribute.Int("statusCode", rw.statusCode),
			),
		)

		// Increment the request counter with the same attributes
		m.requestCounter.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.String("uri", r.RequestURI),
				attribute.Int("statusCode", rw.statusCode),
			),
		)
	}

	return http.HandlerFunc(fn)
}

// WriteHeader captures the status code and delegates to the wrapped ResponseWriter.
// This method intercepts the status code being written to the HTTP response so that
// it can be included in metrics, while maintaining the original functionality.
//
// Parameters:
//   - code: The HTTP status code to write to the response.
func (lrw *responseWriter) WriteHeader(code int) {
	// Store the status code for metrics collection
	lrw.statusCode = code

	// Forward the call to the underlying ResponseWriter
	lrw.ResponseWriter.WriteHeader(code)
}
