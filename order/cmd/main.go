package main

import (
	"context"
	"os"

	"github.com/nico-phil/grpc-microservices/order/config"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/db"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/grpc"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/payment"
	"github.com/nico-phil/grpc-microservices/order/internal/application/core/api"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "order"
	environment = "dev"
	id          = 1
)

func tracerProvider(ctx context.Context) (*tracesdk.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
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

type CustomLogger struct {
	formatter log.JSONFormatter
}

func(l CustomLogger) Format(entry *log.Entry)([]byte, error){
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["space_id"] = span.SpanContext().SpanID().String()
	entry.Data["context"] = span.SpanContext()
	return l.formatter.Format(entry)
}

func init() {
	log.SetFormatter(CustomLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		}},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}


func main(){
	ctx := context.Background()
	tp, err := tracerProvider(ctx)
	if err != nil {
		log.WithContext(ctx).Error("error initializing trace:", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	
	defer tp.Shutdown(ctx)

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