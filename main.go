package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer() *sdktrace.TracerProvider {
	fmt.Println("Initializing Jaeger exporter...")
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://127.0.0.1:14268/api/traces")))
	if err != nil {
		log.Fatalf("failed to create Jaeger exporter: %v", err)
	}
	fmt.Println("Jaeger exporter initialized successfully")

	resource, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("user-service"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	tp := InitTracer()
	defer func() {
		fmt.Println("Shutting down tracer provider...")
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	ctx := context.Background()
	tracer := otel.Tracer("main-tracer")
	ctx, parentSpan := tracer.Start(ctx, "main-function")
	defer parentSpan.End()

	parentSpan.SetAttributes(attribute.String("operation.type", "main"))

	fmt.Println("Starting main operation...")
	time.Sleep(100 * time.Millisecond)

	CallChildFunction(ctx)

	fmt.Println("Continuing main operation...")
	time.Sleep(100 * time.Millisecond)

	fmt.Println("All operations completed.")
}

func CallChildFunction(ctx context.Context) {
	tracer := otel.Tracer("child-1-tracer")
	ctx, span := tracer.Start(ctx, "child-1-function")
	defer span.End()

	span.SetAttributes(attribute.String("operation.type", "child-1"))

	fmt.Println("Performing child 1 operation...")
	time.Sleep(200 * time.Millisecond)

	CallChild2Function(ctx)
}

func CallChild2Function(ctx context.Context) {
	tracer := otel.Tracer("child-2-tracer")
	_, span := tracer.Start(ctx, "child-2-function")
	defer span.End()

	span.SetAttributes(attribute.String("operation.type", "child-2"))

	fmt.Println("Performing child 2 operation...")
	time.Sleep(150 * time.Millisecond)
}
