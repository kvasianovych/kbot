/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// TelemetryShutdown holds cleanup functions for telemetry
type TelemetryShutdown struct {
	tracerProvider *trace.TracerProvider
	meterProvider  *metric.MeterProvider
}

// Shutdown gracefully shuts down telemetry providers
func (t *TelemetryShutdown) Shutdown(ctx context.Context) {
	if t.tracerProvider != nil {
		if err := t.tracerProvider.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown tracer provider", "error", err)
		}
	}
	if t.meterProvider != nil {
		if err := t.meterProvider.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown meter provider", "error", err)
		}
	}
}

// InitTelemetry initializes OpenTelemetry tracing and metrics
// It reads configuration from environment variables:
// - OTEL_EXPORTER_OTLP_ENDPOINT: collector endpoint
// - OTEL_SERVICE_NAME: service name (default: kbot)
func InitTelemetry(ctx context.Context) (*TelemetryShutdown, error) {
	shutdown := &TelemetryShutdown{}

	// Get service name from env or use default
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "kbot"
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(appVersion),
		),
		resource.WithFromEnv(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, err
	}

	// Initialize tracer provider
	tracerProvider, err := initTracerProvider(ctx, res)
	if err != nil {
		return nil, err
	}
	shutdown.tracerProvider = tracerProvider
	otel.SetTracerProvider(tracerProvider)

	// Set up propagator for distributed tracing
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Initialize meter provider
	meterProvider, err := initMeterProvider(ctx, res)
	if err != nil {
		return nil, err
	}
	shutdown.meterProvider = meterProvider
	otel.SetMeterProvider(meterProvider)

	slog.Info("telemetry initialized", "service", serviceName, "version", appVersion)
	return shutdown, nil
}

func initTracerProvider(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	return tp, nil
}

func initMeterProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	return mp, nil
}
