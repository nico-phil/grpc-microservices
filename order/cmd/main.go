package main

import (
	"log"

	"github.com/nico-phil/grpc-microservices/order/config"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/db"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/grpc"
	"github.com/nico-phil/grpc-microservices/order/internal/adapters/payment"
	"github.com/nico-phil/grpc-microservices/order/internal/application/core/api"
)

func main(){
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