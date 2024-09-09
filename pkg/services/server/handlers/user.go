package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/marketplace/product"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/session"
	"mephiMainProject/pkg/services/server/user"
	"net/http"
	"time"
)

type UserHandler struct {
	Logger   *zap.SugaredLogger
	Sessions *session.SessionManager
	UserRepo user.UserRepo
}

type FormData struct {
	Values   map[string]string
	Errors   map[string]string
	Products []*product.Product
}

func NewFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	FormData FormData
}

func (h *UserHandler) LoginGET(c echo.Context) error {
	formData := NewFormData()
	formData.Values["requestType"] = "login"
	return c.Render(http.StatusOK, "index", formData)
}

func (h *UserHandler) RegisterGET(c echo.Context) error {
	formData := NewFormData()
	formData.Values["requestType"] = "register"
	return c.Render(http.StatusOK, "index", formData)
}

func (h *UserHandler) LoginPOST(c echo.Context) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	if login == "" || password == "" {
		formData := NewFormData()
		formData.Values["requestType"] = "login"
		formData.Values["login"] = login
		formData.Values["password"] = password
		formData.Errors["error"] = "Login or password is not correct"
		return c.Render(422, "form", formData)
	}
	userAuthData, err := h.UserRepo.Authorize(login, password)
	if err != nil {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "login"
		formData.Values["password"] = password
		formData.Errors["error"] = "Invalid login or password"
		return c.Render(422, "form", formData)
	}
	newSession, err := h.Sessions.Create(userAuthData.Login)
	if err != nil {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "login"
		formData.Errors["error"] = err.Error()
		return c.Render(422, "form", formData)
	}

	h.Logger.Infof("Successfully created session for username %v", newSession.Username)
	token, err := session.CreateNewToken(config.User{Login: userAuthData.Login}, newSession.SessID.ID)
	if err != nil {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "login"
		formData.Errors["error"] = err.Error()
		return c.Render(422, "form", formData)
	}
	c.SetCookie(&http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Second * 60 * 60 * 24 * 3),
	})
	h.Logger.Infof("Send token on client for user with username: %v ", newSession.Username)
	return c.Redirect(http.StatusSeeOther, "/")
}

// TODO: сделать отправку подзапроса на проверку занятости никнейма
func (h *UserHandler) RegisterPOST(c echo.Context) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	if login == "" || password == "" {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "register"
		formData.Values["password"] = password
		formData.Errors["error"] = "Try other"
		return c.Render(422, "form", formData)
	}
	_, err := h.UserRepo.Register(login, password)
	if err != nil {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "register"
		formData.Values["password"] = password
		formData.Errors["error"] = err.Error()
		return c.Render(422, "form", formData)
	}
	newSession, err := h.Sessions.Create(login)
	if err != nil {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "register"
		formData.Errors["error"] = err.Error()
		return c.Render(422, "form", formData)
	}

	h.Logger.Infof("Successfully created session for username %v", newSession.Username)
	token, err := session.CreateNewToken(config.User{Login: login}, newSession.SessID.ID)
	if err != nil {
		formData := NewFormData()
		formData.Values["login"] = login
		formData.Values["requestType"] = "register"
		formData.Errors["error"] = err.Error()
		return c.Render(422, "form", formData)
	}
	c.SetCookie(&http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Second * 60 * 60 * 24 * 3),
	})
	h.Logger.Infof("Send token on client for user with username: %v ", newSession.Username)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *UserHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "session",
		Value:   "expired",
		Expires: time.Now().Add(-1 * time.Hour * 24),
	},
	)
	return c.Redirect(http.StatusSeeOther, "/")
}
