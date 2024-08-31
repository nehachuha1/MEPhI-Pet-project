package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/profile"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
)

type ProfileHandler struct {
	Logger      *zap.SugaredLogger
	Sessions    *session.SessionManager
	ProfileRepo profile.ProfileRepo
}

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant read request body")
	}

	userToCreate := &config.User{}
	err = json.Unmarshal(body, userToCreate)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	currentSession, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Auth required", http.StatusUnauthorized)
		return
	}
	err = h.ProfileRepo.CreateProfile(userToCreate, currentSession.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.Logger.Infof("Successfull created profile for user %v", currentSession.Username)
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}
	vars := mux.Vars(r)
	currentUsername, isOk := vars["USERNAME"]
	if !isOk {
		http.Error(w, "Current URl hasn't username", http.StatusBadRequest)
		return
	}
	user, err := h.ProfileRepo.GetProfile(currentUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(user)
	CheckMarshalError(w, err, resp)
	w.WriteHeader(http.StatusOK)
	h.Logger.Infof("Send profile of %v", currentUsername)
}

func (h *ProfileHandler) EditProfile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant read request body")
	}
	newUserData := &config.User{}
	err = json.Unmarshal(body, newUserData)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	currentSession, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Auth required", http.StatusUnauthorized)
		return
	}
	err = h.ProfileRepo.EditProfile(currentSession.Username, newUserData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	h.Logger.Infof("Successfull edited profile of user %v", currentSession.Username)
}

func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	currentSession, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Auth required", http.StatusUnauthorized)
		return
	}
	err = h.ProfileRepo.DeleteProfile(currentSession.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.Logger.Infof("Successfull delete profile of user %v", currentSession.Username)
}
