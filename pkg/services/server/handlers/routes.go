package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/ownMiddleware"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
)

func GenerateRoutes(currentCfg *config.Config, sm *session.SessionManager, uh UserHandler, mh MarketplaceHandler, ph ProfileHandler) *echo.Echo {
	e := echo.New()

	//static connection
	e.Renderer = config.NewTemplates()
	e.Static("/views/css", "views/css")

	e.GET("/", func(c echo.Context) error {
		formData := NewFormData()
		currentSession, err := session.SessionFromContext(c)
		if err != nil {
			log.Printf("MAIN PAGE ERR - %v\n", err)
			return c.Render(http.StatusOK, "index", formData)
		}
		formData.Values["username"] = currentSession.Username
		return c.Render(http.StatusOK, "index", formData)
	})

	// authorization handlers
	e.GET("/login", uh.LoginGET)
	e.GET("/register", uh.RegisterGET)
	e.GET("/logout", uh.Logout)
	e.POST("/register", uh.RegisterPOST)
	e.POST("/login", uh.LoginPOST)

	e.GET("/profile/:username", ph.GetProfile)
	e.POST("/profile/create", ph.CreateProfile)

	//marketplace handlers
	e.GET("/marketplace", mh.GetProducts)

	//not found
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.String(http.StatusOK, "Under construct")
	})

	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ownMiddleware.Auth(sm))

	return e
}
