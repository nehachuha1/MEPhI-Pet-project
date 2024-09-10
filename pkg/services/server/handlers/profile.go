package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/profile"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
	"strconv"
	"time"
)

type ProfileHandler struct {
	Logger      *zap.SugaredLogger
	ProfileRepo profile.ProfileRepo
}

func (h *ProfileHandler) CreateProfile(c echo.Context) error {
	currentSession, _ := session.SessionFromContext(c)
	currentProfile := &config.User{
		Login:        currentSession.Username,
		FirstName:    c.FormValue("name"),
		SecondName:   c.FormValue("surname"),
		Sex:          c.FormValue("sex"),
		Address:      c.FormValue("address"),
		RegisterDate: time.Now().Format("01-02-2006 15:04:05"),
		EditDate:     time.Now().Format("01-02-2006 15:04:05"),
	}
	currentProfile.Age, _ = strconv.Atoi(c.FormValue("age"))
	err := h.ProfileRepo.CreateProfile(currentProfile, currentSession.Username)
	formData := NewFormData()
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "profile", formData)
	}
	return c.Redirect(http.StatusSeeOther, "/profile/"+currentSession.Username)
}

func (h *ProfileHandler) GetProfile(c echo.Context) error {
	formData := NewFormData()
	currentProfile, err := h.ProfileRepo.GetProfile(c.Param("username"))
	if err != nil {
		return c.Render(200, "profile", formData)
	}
	formData.Profile["name"] = currentProfile.FirstName
	formData.Profile["surname"] = currentProfile.SecondName
	formData.Profile["sex"] = currentProfile.Sex
	formData.Profile["age"] = currentProfile.Age
	formData.Profile["address"] = currentProfile.Address
	formData.Profile["registerDate"] = currentProfile.RegisterDate
	formData.Profile["editDate"] = currentProfile.EditDate
	return c.Render(200, "profile", formData)
}
