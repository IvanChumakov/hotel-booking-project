package tracing

import (
	"context"
	"google.golang.org/grpc/metadata"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewTrace() error {
	tracingUrl, _ := os.LookupEnv("TRACING_URL")

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(tracingUrl)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("hotel-booking-service"),
		)),
	)

	otel.SetTracerProvider(tp)
	return nil
}

func StartTracerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("hotel-service").Start(ctx, "hotel-booking-service:"+spanName)
}

func GetParentContext(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceIdString := md["x-trace-id"][0]

	traceId, err := trace.TraceIDFromHex(traceIdString)
	if err != nil {
		return nil, err
	}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})
	return trace.ContextWithSpanContext(ctx, spanContext), nil
}

func GetParentContextFromHeader(ctx context.Context, traceIdStr string) (context.Context, error) {
	traceId, err := trace.TraceIDFromHex(traceIdStr)
	if err != nil {
		return nil, err
	}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})
	return trace.ContextWithSpanContext(ctx, spanContext), nil
}
