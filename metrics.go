// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package metrics provides tools for collecting, exporting, and monitoring metrics
// in Go applications with support for OpenTelemetry Protocol (OTLP).
//
// The package offers integrations with various exporters and custom metric collectors,
// allowing developers to easily implement application monitoring. It's designed to
// work seamlessly with other GoKit packages like configs and logging.
package metrics

import (
	"github.com/goxkit/configs"
	"github.com/goxkit/metrics/noop"
	"github.com/goxkit/metrics/otlp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Install initializes and configures a metric provider based on the application's configuration.
// It determines whether to use the OpenTelemetry Protocol (OTLP) exporter or a no-operation
// implementation depending on the configuration.
//
// Parameters:
//   - cfgs: Application configuration containing metrics settings
//
// Returns:
//   - A configured OpenTelemetry MeterProvider
//   - An error if the initialization fails
func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error) {
	if cfgs.OTLPConfigs.Enabled {
		return otlp.Install(cfgs)
	}

	return noop.Install(cfgs)
}
