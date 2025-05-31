// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

package metrics

import (
	"github.com/goxkit/configs"
	"github.com/goxkit/metrics/noop"
	"github.com/goxkit/metrics/otlp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error) {
	if cfgs.OTLPConfigs.Enabled {
		return otlp.Install(cfgs)
	}

	return noop.Install(cfgs)
}
