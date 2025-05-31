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

func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()

	if cfgs.OTLPExporterConn == nil {
		conn, err := otlpgrpc.NewExporterGRPCClient(cfgs)
		if err != nil {
			cfgs.Logger.Error("failed to create grpc exporter", zap.Error(err))
			return nil, err
		}
		cfgs.OTLPExporterConn = conn
	}

	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(cfgs.OTLPExporterConn),
	)
	if err != nil {
		cfgs.Logger.Error("failed to create OTLP trace exporter", zap.Error(err))
		return nil, err
	}

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

	cfgs.MetricsProvider = meterProvider
	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}
