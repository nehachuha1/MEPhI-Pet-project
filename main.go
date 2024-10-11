package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"mephiMainProject/pkg/services/marketplace/orders"
	"mephiMainProject/pkg/services/marketplace/product"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/database"
	"mephiMainProject/pkg/services/server/handlers"
	"mephiMainProject/pkg/services/server/profile"
	"mephiMainProject/pkg/services/server/session"
	"mephiMainProject/pkg/services/server/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Can't load .env file")
	}

	currentCfg := config.NewConfig()
	dbControl := database.NewDBUsage(currentCfg)
	userRepo := user.NewUserRepository(currentCfg)
	profileRepo := profile.NewProfileRepository(currentCfg)
	sessionManager := session.NewSessionManager(dbControl, currentCfg)

	//	GRPC connection to services
	grpcConn, err := grpc.NewClient(":"+currentCfg.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC connection err - %v", err)
		return
	}
	defer grpcConn.Close()

	if err != nil {
		log.Fatalf("gRPC starting err - %v", err)
		return
	}
	marketplaceService := product.NewMarketplaceServiceClient(grpcConn)
	orderService := orders.NewOrderServiceClient(grpcConn)

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err = zapLogger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()
	logger := zapLogger.Sugar()

	userHandler := handlers.UserHandler{
		Logger:   logger,
		Sessions: sessionManager,
		UserRepo: userRepo,
	}

	profileHandler := handlers.ProfileHandler{
		Logger:      logger,
		ProfileRepo: profileRepo,
	}

	marketHandler := handlers.MarketplaceHandler{
		Logger:             logger,
		MarketPlaceManager: marketplaceService,
		OrdersManager:      orderService,
		ProfileRepo:        profileRepo,
	}
	//
	echoHandler := handlers.GenerateRoutes(currentCfg, sessionManager, userHandler, marketHandler, profileHandler)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	echoHandler.Logger.Fatal(echoHandler.Start(addr))
}
