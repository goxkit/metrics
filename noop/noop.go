// Package noop provides a no-operation implementation of the metrics system.
// This is useful for development, testing, or when metrics collection is not needed,
// as it implements the expected interfaces without actually collecting or exporting any data.
package noop

import (
	"github.com/goxkit/configs"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Install creates and configures a no-operation metrics provider.
// It creates an empty MeterProvider that doesn't perform any actual metrics collection
// and stores it in the application configuration for later use.
//
// Parameters:
//   - cfgs: Application configuration where the metrics provider will be stored
//
// Returns:
//   - A configured no-operation MeterProvider that satisfies the interface requirements
//   - Always returns nil error since this implementation cannot fail
func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error) {
	provider := sdkmetric.NewMeterProvider()
	cfgs.MetricsProvider = provider
	return provider, nil
}
