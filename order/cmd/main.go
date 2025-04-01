package main

import (
	"context"
	"log"

	"github.com/nico-phil/grpc-microservices/order/config"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/db"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/grpc"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/payment"
	"github.com/nico-phil/grpc-microservices/order/internal/application/core/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "order"
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

	// tp, err := tracerProvider()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// otel.SetTracerProvider(tp)
	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	dbAdapter, err := db.NewAdapter(config.GetDataSourceUrl())
	if err != nil {
		log.Fatalf("failed to connect to db. Error %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("failed to initialize payment stub. error %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}