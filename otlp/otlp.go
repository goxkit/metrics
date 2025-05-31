// Package otlp provides an implementation of the metrics system using the OpenTelemetry Protocol.
// It configures and sets up a metrics exporter that sends data to an OTLP-compatible collector
// using gRPC transport.
package otlp

import (
	"context"

	"github.com/goxkit/configs"
	"github.com/goxkit/otel/otlpgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

// Install creates and configures an OpenTelemetry Protocol (OTLP) metrics provider.
// It sets up a gRPC connection to the configured OTLP endpoint, creates an exporter,
// and initializes a MeterProvider with appropriate resource attributes.
//
// Parameters:
//   - cfgs: Application configuration containing OTLP settings and where the metrics provider will be stored
//
// Returns:
//   - A configured MeterProvider that exports metrics via OTLP
//   - An error if any part of the configuration process fails
func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()

	// Create a gRPC client connection if one doesn't exist yet
	if cfgs.OTLPExporterConn == nil {
		conn, err := otlpgrpc.NewExporterGRPCClient(cfgs)
		if err != nil {
			cfgs.Logger.Error("failed to create grpc exporter", zap.Error(err))
			return nil, err
		}
		cfgs.OTLPExporterConn = conn
	}

	// Create the OTLP metrics exporter using the gRPC connection
	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(cfgs.OTLPExporterConn),
	)
	if err != nil {
		cfgs.Logger.Error("failed to create OTLP metric exporter", zap.Error(err))
		return nil, err
	}

	// Create the meter provider with periodic collection and resource attributes
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfgs.AppConfigs.Name),
			semconv.ServiceNamespaceKey.String(cfgs.AppConfigs.Namespace),
			attribute.String("service.environment", cfgs.AppConfigs.Environment.String()),
			semconv.DeploymentEnvironmentKey.String(cfgs.AppConfigs.Environment.String()),
			semconv.TelemetrySDKLanguageKey.String("go"),
			semconv.TelemetrySDKLanguageGo.Key.Bool(true),
		)),
	)

	// Store the provider in the configs and set as global provider
	cfgs.MetricsProvider = meterProvider
	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}
