package handlers

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/server/middleware"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
)

func GenerateRoutes(uh UserHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", uh.Login).Methods("POST")
	r.HandleFunc("/register", uh.Register).Methods("POST")

	return r
}

func AddProcessing(r *mux.Router, sm *session.SessionManager, logger *zap.SugaredLogger) http.Handler {
	r.Use(middleware.Auth(sm))
	r.Use(middleware.AccessLog(logger))
	r.Use(middleware.Panic)
	return r
}
