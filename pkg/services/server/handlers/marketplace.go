package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"mephiMainProject/pkg/services/marketplace/product"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
	"strconv"
)

type MarketplaceHandler struct {
	Logger             *zap.SugaredLogger
	Sessions           *session.SessionManager
	MarketPlaceManager product.MarketplaceServiceClient
}

func (mh *MarketplaceHandler) ListAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}

	allProducts, err := mh.MarketPlaceManager.GetAllProducts(
		r.Context(), &product.Nothing{Dummy: true},
	)
	if err != nil {
		http.Error(w, "ListAllProducts err - "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(allProducts)
	w.WriteHeader(http.StatusOK)
	CheckMarshalError(w, err, resp)
	mh.Logger.Infoln("Successfully listed all products from DB")
}

func (mh *MarketplaceHandler) ListProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}
	requestVars := mux.Vars(r)
	productID := requestVars["PRODUCT_ID"]

	productFromDB, err := mh.MarketPlaceManager.GetProduct(
		r.Context(), &product.ProductID{ProductID: productID},
	)

	if err != nil {
		http.Error(w, "List product handler - "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(productFromDB)
	w.WriteHeader(http.StatusOK)
	CheckMarshalError(w, err, resp)
	mh.Logger.Infof("Successfully got product with id %v from DB\n", productID)
}

func (mh *MarketplaceHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}

	_, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Auth is required", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newProduct := &product.Product{}
	err = json.Unmarshal(body, newProduct)
	if err != nil {
		http.Error(w, "Unmarshalling error", http.StatusBadRequest)
	}

	resp, err := mh.MarketPlaceManager.CreateProduct(r.Context(), newProduct)
	returnResponse, err := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	CheckMarshalError(w, err, returnResponse)
	mh.Logger.Infof("Successfully created product with name %s from DB\n", newProduct.Name)
}

func (mh *MarketplaceHandler) EditProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}

	ownerUsername, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Auth is required", http.StatusUnauthorized)
		return
	}

	requestVars := mux.Vars(r)
	productID, err := strconv.ParseUint(requestVars["PRODUCT_ID"], 10, 64) // returns uint64
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	editedData := &product.Product{}
	err = json.Unmarshal(body, editedData)
	if err != nil {
		http.Error(w, "Unmarshalling error", http.StatusBadRequest)
	}
	editedData.Id = int64(productID)
	editedData.OwnerUsername = ownerUsername.Username

	resp, err := mh.MarketPlaceManager.EditProduct(r.Context(), editedData)
	returnResponse, err := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	CheckMarshalError(w, err, returnResponse)
	mh.Logger.Infof("Successfully edited product with name %s from DB\n", editedData.Name)
}

func (mh *MarketplaceHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}

	_, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Auth is required", http.StatusUnauthorized)
		return
	}

	requestVars := mux.Vars(r)
	productID := requestVars["PRODUCT_ID"]

	pID := &product.ProductID{ProductID: productID}

	resp, err := mh.MarketPlaceManager.DeleteProduct(r.Context(), pID)
	returnResponse, err := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	CheckMarshalError(w, err, returnResponse)
	mh.Logger.Infof("Successfully deleted product with id %v from DB\n", pID.ProductID)
}

// list/list all <- get request
// create/edit <- post request
// delete <- delete request
