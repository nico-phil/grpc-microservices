package main

import (
	"log"

	"github.com/nico-phil/grpc-microservices/payment/config"
	"github.com/nico-phil/grpc-microservices/payment/internal/adapters/db"
	"github.com/nico-phil/grpc-microservices/payment/internal/adapters/grpc"
	"github.com/nico-phil/grpc-microservices/payment/internal/application/core/api"
)

func main(){
	dbApater, err := db.NewAdapter(config.GetDataSourceUrl())
	if err != nil {
		log.Fatalf("failed to connect to db error: %v", err)
	}
	
	application  := api.NewApplication(dbApater)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())

	grpcAdapter.Run()
}