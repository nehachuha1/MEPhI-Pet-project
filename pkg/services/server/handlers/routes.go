package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/ownMiddleware"
	"mephiMainProject/pkg/services/server/session"
)

func GenerateRoutes(sm *session.SessionManager, uh UserHandler) *echo.Echo {
	e := echo.New()

	//static connection
	e.Renderer = config.NewTemplates()
	e.Static("/views/css", "views/css")

	// authorization handlers
	e.GET("/login", uh.LoginGET)
	e.GET("/register", uh.RegisterGET)
	e.POST("/register", uh.RegisterPOST)
	e.POST("/login", uh.LoginPOST)

	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ownMiddleware.Auth(sm))

	return e
}
