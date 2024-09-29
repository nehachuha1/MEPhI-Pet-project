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
	formData := NewFormData()
	currentProfile := &config.User{
		Login:        currentSession.Username,
		FirstName:    c.FormValue("name"),
		SecondName:   c.FormValue("surname"),
		RegisterDate: time.Now().Format("01-02-2006 15:04:05"),
		EditDate:     time.Now().Format("01-02-2006 15:04:05"),
	}
	age, _ := strconv.Atoi(c.FormValue("age"))
	if age < 0 || age > 100 {
		formData.Errors["createError"] = "Invalid age"
	}
	sex := c.FormValue("sex")
	if sex == "F" || sex == "M" {
		currentProfile.Sex = sex
	} else {
		formData.Errors["createError"] = "Invalid sex"
	}
	currentProfile.Address = c.FormValue("address") + " | Room: " + c.FormValue("room")
	currentProfile.Age = age
	err := h.ProfileRepo.CreateProfile(currentProfile, currentSession.Username)
	formData.Values["username"] = currentSession.Username
	if currentSession.Username == currentProfile.Login {
		formData.Values["currentUserIsOwner"] = currentSession.Username
	}
	if err != nil {
		formData.Errors["createError"] = err.Error()
		return c.Render(422, "profile", formData)
	}
	formData.Profile["name"] = currentProfile.FirstName
	formData.Profile["surname"] = currentProfile.SecondName
	formData.Profile["sex"] = currentProfile.Sex
	formData.Profile["age"] = currentProfile.Age
	formData.Profile["address"] = currentProfile.Address
	formData.Profile["registerDate"] = currentProfile.RegisterDate
	formData.Profile["editDate"] = currentProfile.EditDate
	c.Render(200, "profile", formData)
	return c.Redirect(http.StatusSeeOther, "/profile/"+currentSession.Username)
}

func (h *ProfileHandler) GetProfile(c echo.Context) error {
	formData := NewFormData()
	sess, _ := session.SessionFromContext(c)
	formData.Values["username"] = sess.Username
	currentProfile, err := h.ProfileRepo.GetProfile(c.Param("username"))
	if err != nil && currentProfile.Login != sess.Username {
		formData.Errors["error"] = err.Error()
		return c.Render(200, "profile-view", formData)
	} else if err != nil && currentProfile.Login == sess.Username {
		formData.Values["currentUserIsOwner"] = sess.Username
		return c.Render(200, "profile-view", formData)
	}
	if sess.Username == currentProfile.Login {
		formData.Values["currentUserIsOwner"] = sess.Username
	}
	formData.Profile["name"] = currentProfile.FirstName
	formData.Profile["surname"] = currentProfile.SecondName
	formData.Profile["sex"] = currentProfile.Sex
	formData.Profile["age"] = currentProfile.Age
	formData.Profile["address"] = currentProfile.Address
	formData.Profile["registerDate"] = currentProfile.RegisterDate
	formData.Profile["editDate"] = currentProfile.EditDate
	return c.Render(200, "profile-view", formData)
}

func (h *ProfileHandler) EditProfileGET(c echo.Context) error {
	currentSession, _ := session.SessionFromContext(c)
	currentProfile, err := h.ProfileRepo.GetProfile(currentSession.Username)
	formData := NewFormData()
	formData.Values["username"] = currentSession.Username
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "profile-edit", formData)
	}
	formData.Values["name"] = currentProfile.FirstName
	formData.Values["surname"] = currentProfile.SecondName
	formData.Values["sex"] = currentProfile.Sex
	formData.Values["age"] = strconv.Itoa(currentProfile.Age)
	formData.Values["address"] = currentProfile.Address
	return c.Render(200, "profile-edit", formData)
}

func (h *ProfileHandler) EditProfilePOST(c echo.Context) error {
	currentSession, _ := session.SessionFromContext(c)
	formData := NewFormData()
	formData.Values["username"] = currentSession.Username
	oldData, _ := h.ProfileRepo.GetProfile(currentSession.Username)
	newData := &config.User{
		Login:        currentSession.Username,
		FirstName:    c.FormValue("name"),
		SecondName:   c.FormValue("surname"),
		RegisterDate: oldData.RegisterDate,
		EditDate:     time.Now().Format("01-02-2006 15:04:05"),
	}
	sex := c.FormValue("sex")
	if sex == "F" || sex == "M" {
		newData.Sex = sex
	} else {
		formData.Errors["error"] = "Invalid sex"
		return c.Render(422, "profile-edit", formData)
	}
	age, _ := strconv.Atoi(c.FormValue("age"))
	if age < 0 || age > 100 {
		formData.Errors["error"] = "Invalid age"
		return c.Render(422, "profile-edit", formData)
	}
	newData.Age = age
	room := c.FormValue("room")
	_, err := strconv.Atoi(room)
	if err != nil {
		formData.Errors["error"] = "Invalid room"
		return c.Render(422, "profile-edit", formData)
	}
	newData.Address = c.FormValue("address") + " | Room: " + room
	err = h.ProfileRepo.EditProfile(currentSession.Username, newData)
	if err != nil {
		formData.Values["name"] = oldData.FirstName
		formData.Values["surname"] = oldData.SecondName
		formData.Values["sex"] = oldData.Sex
		formData.Values["age"] = strconv.Itoa(oldData.Age)
		formData.Values["address"] = oldData.Address

		formData.Errors["error"] = err.Error()
		return c.Render(422, "profile-edit", formData)
	}
	currentProfile, err := h.ProfileRepo.GetProfile(currentSession.Username)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(200, "profile-view", formData)
	}
	if currentSession.Username == currentProfile.Login {
		formData.Values["currentUserIsOwner"] = currentSession.Username
	}
	formData.Profile["name"] = currentProfile.FirstName
	formData.Profile["surname"] = currentProfile.SecondName
	formData.Profile["sex"] = currentProfile.Sex
	formData.Profile["age"] = currentProfile.Age
	formData.Profile["address"] = currentProfile.Address
	formData.Profile["registerDate"] = currentProfile.RegisterDate
	formData.Profile["editDate"] = currentProfile.EditDate
	c.Render(200, "profile", formData)
	return c.Redirect(http.StatusSeeOther, "/profile/"+currentSession.Username)
}
