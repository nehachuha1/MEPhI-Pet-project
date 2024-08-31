package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/database"
	"mephiMainProject/pkg/services/server/handlers"
	"mephiMainProject/pkg/services/server/profile"
	"mephiMainProject/pkg/services/server/session"
	"mephiMainProject/pkg/services/server/user"
	"net/http"
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

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := zapLogger.Sync()
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
		Sessions:    sessionManager,
		ProfileRepo: profileRepo,
	}

	addHandlersMux := handlers.GenerateRoutes(userHandler, profileHandler)
	addProcessing := handlers.AddProcessing(addHandlersMux, sessionManager, logger)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, addProcessing)
	if err != nil {
		fmt.Println(err)
	}
}
