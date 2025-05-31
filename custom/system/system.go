// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package system provides system metrics collection capabilities for monitoring
// memory usage, garbage collection, threads, and goroutines in Go applications.
package system

import (
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// BasicMetricsCollector initializes and configures basic system metrics collection.
// It sets up memory and system gauges and starts the continuous collection of metrics
// to monitor runtime performance and resource usage of the application.
//
// Parameters:
//   - logger: A logger instance for logging metrics-related messages.
//
// Returns:
//   - An error if metrics collection could not be initialized.
func BasicMetricsCollector(logger *zap.SugaredLogger) error {
	logger.Debug("configuring basic metrics...")

	// Create a meter with an appropriate instrumentation scope name
	meter := otel.Meter("github.com/goxkit/metrics/custom/system")

	// Initialize memory statistics collection
	mem, err := NewMemGauges(meter)
	if err != nil {
		return err
	}

	// Initialize system statistics collection (threads, goroutines, etc.)
	sys, err := NewSysGauge(meter)
	if err != nil {
		return err
	}

	logger.Debug("basic metrics configured")

	// Start collecting metrics by registering the callbacks
	mem.Collect(meter)
	sys.Collect(meter)

	return nil
}
