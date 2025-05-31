// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package system provides system-level metrics collection for monitoring
// runtime characteristics of Go applications.
package system

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

// NewSysGauge creates a new system metrics collector that monitors
// OS threads, CGO calls, and active goroutines. These metrics provide
// insights into the concurrency patterns and resource utilization of
// the Go application.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to create gauge instruments.
//
// Returns:
//   - A BasicGauges implementation for system metrics collection.
//   - An error if any gauge creation fails.
func NewSysGauge(meter metric.Meter) (BasicGauges, error) {
	// Create a gauge for tracking the number of OS threads
	ggThreads, err := meter.Int64ObservableGauge("go_threads", metric.WithDescription("Number of OS threads created."))
	if err != nil {
		return nil, err
	}

	// Create a gauge for tracking the number of CGO calls
	ggCgo, err := meter.Int64ObservableGauge("go_cgo", metric.WithDescription("Number of CGO calls."))
	if err != nil {
		return nil, err
	}

	// Create a gauge for tracking the number of goroutines
	ggGRoutines, err := meter.Int64ObservableGauge("go_goroutines", metric.WithDescription("Number of goroutines."))
	if err != nil {
		return nil, err
	}

	// Return the configured system gauges
	return &sysGauges{
		ggThreads, ggCgo, ggGRoutines,
	}, nil
}

// Collect registers callbacks for system metrics collection.
// It reads statistics from the Go runtime about CPU cores, CGO calls,
// and goroutines and reports them through the observable gauges.
// This provides visibility into the application's concurrency behavior
// and resource utilization.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to register callbacks.
func (s *sysGauges) Collect(meter metric.Meter) {
	// Define the callback function for collecting system metrics
	cb := func(_ context.Context, observer metric.Observer) error {
		// Record the number of CPU cores available
		observer.ObserveInt64(s.ggThreads, int64(runtime.NumCPU()))

		// Record the number of CGO calls made
		observer.ObserveInt64(s.ggCgo, int64(runtime.NumCgoCall()))

		// Record the number of currently active goroutines
		observer.ObserveInt64(s.ggGRoutines, int64(runtime.NumGoroutine()))

		return nil
	}

	// Register the callback with the meter
	_, _ = meter.RegisterCallback(cb)
}
