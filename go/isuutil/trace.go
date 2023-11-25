package isuutil

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	// traceIDRatioBased はトレースのサンプリングレート
	// Jaegerが一度に表示できるトレース数が1500なので、パフォーマンス改善に応じて調整する
	traceIDRatioBased = 0.01
	// 基本的にWatchTowerにホストされているJaegerにTraceを送る
	endpoint = "watchtower:4317"

	serviceName = "webapp"
)

// InitializeTracerProvider はmain関数で呼び出されることを想定。
// 設定したTracerは otel.GetTracerProvider().Tracer("") で呼び出し可能。
func InitializeTracerProvider() (*sdktrace.TracerProvider, error) {
	res, err := resource.New(context.Background(),
		resource.WithTelemetrySDK(),
	)
	resAttr := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	if err != nil {
		return nil, fmt.Errorf("faield to create resource: %w", err)
	}

	client := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	exp, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	tracerProviderOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(traceIDRatioBased))),
		sdktrace.WithResource(res),
		sdktrace.WithResource(resAttr),
		sdktrace.WithBatcher(exp),
	}

	tp := sdktrace.NewTracerProvider(tracerProviderOptions...)
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	return tp, nil
}
