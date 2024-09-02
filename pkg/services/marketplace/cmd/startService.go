package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"mephiMainProject/pkg/services/marketplace/config"
	"mephiMainProject/pkg/services/marketplace/product"
	"net"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Can't load .env file")
	}

	currentConfig := config.NewConfig()

	marketPlaceServ := product.NewMarketplaceService(currentConfig)

	listener, err := net.Listen("tcp", ":"+currentConfig.GRPC.Port)
	if err != nil {
		log.Fatalf("Can't start marketp[lace service. Err - %v\n", err)
	}

	server := grpc.NewServer()
	product.RegisterMarketplaceServiceServer(server, marketPlaceServ)

	log.Printf("Successfully started marketplace service\n")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
