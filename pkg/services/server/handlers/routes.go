package handlers

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/server/middleware"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
)

func GenerateRoutes(uh UserHandler, ph ProfileHandler, mh MarketplaceHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", uh.Login).Methods("POST")
	r.HandleFunc("/api/register", uh.Register).Methods("POST")

	r.HandleFunc("/api/profile/create", ph.CreateProfile).Methods("POST")
	r.HandleFunc("/api/profile/{USERNAME}", ph.GetProfile).Methods("GET")
	r.HandleFunc("/api/profile/edit", ph.EditProfile).Methods("POST")
	r.HandleFunc("/api/profile/delete", ph.DeleteProfile).Methods("POST")

	r.HandleFunc("/api/marketplace/products", mh.ListAllProducts).Methods("GET")
	r.HandleFunc("/api/marketplace/products/{PRODUCT_ID}", mh.ListProduct).Methods("GET")
	r.HandleFunc("/api/marketplace/create", mh.CreateProduct).Methods("POST")
	r.HandleFunc("/api/marketplace/{PRODUCT_ID}/edit", mh.EditProduct).Methods("POST")
	r.HandleFunc("/api/marketplace/{PRODUCT_ID}/delete", mh.DeleteProduct).Methods("DELETE")

	return r
}

func AddProcessing(r *mux.Router, sm *session.SessionManager, logger *zap.SugaredLogger) http.Handler {
	r.Use(middleware.Auth(sm))
	r.Use(middleware.AccessLog(logger))
	r.Use(middleware.Panic)
	return r
}
