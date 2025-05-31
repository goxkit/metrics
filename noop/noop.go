package noop

import (
	"github.com/goxkit/configs"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error) {
	provider := sdkmetric.NewMeterProvider()
	cfgs.MetricsProvider = provider
	return provider, nil
}
