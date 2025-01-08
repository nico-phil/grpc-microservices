package main

import (
	"context"
	"log"

	"github.com/nico-phil/grpc-microservices/payment/config"
	"github.com/nico-phil/grpc-microservices/payment/internal/adapters/db"
	"github.com/nico-phil/grpc-microservices/payment/internal/adapters/grpc"
	"github.com/nico-phil/grpc-microservices/payment/internal/application/core/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "payment"
	environment = "dev"
	id          = 1
)

func tracerProvider() (*tracesdk.TracerProvider, error) {
	exp, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
		
	)
	return tp, nil
}

func main(){
	tp, err :=tracerProvider()
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	dbApater, err := db.NewAdapter(config.GetDataSourceUrl())
	if err != nil {
		log.Fatalf("failed to connect to db error: %v", err)
	}
	
	application  := api.NewApplication(dbApater)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())

	grpcAdapter.Run()
}