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
	e.Static("/", "public")
	e.Static("/views/css", "views/css")
	e.Static("/data/img", "data/img")
	e.Static("/views/js/", "views/js")
	e.File("/favicon.ico", "views/favicon.ico")

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

	//profile handlers
	e.GET("/profile/:username", ph.GetProfile)
	e.POST("/profile/create", ph.CreateProfile)
	e.GET("/profile/edit", ph.EditProfileGET)
	e.POST("/profile/edit", ph.EditProfilePOST)

	//marketplace handlers
	e.DELETE("/marketplace/product/:id/delete", mh.DeleteProduct)
	e.GET("/marketplace/product/:id", mh.GetProduct)
	e.GET("/marketplace/products/:username", mh.GetUserProducts)
	e.GET("/marketplace", mh.GetProducts)
	e.GET("/marketplace/create", mh.CreateProductGet)
	e.POST("/marketplace/create", mh.CreateProductPost)

	// Orders
	e.GET("/marketplace/orders/", mh.GetOrders)
	e.GET("/marketplace/sales/", mh.GetSales)
	e.POST("/marketplace/sales/order/:id/accept", mh.AcceptOrder)
	e.POST("/marketplace/sales/order/:id/complete", mh.CompleteOrder)
	e.POST("/marketplace/orders/create", mh.ProceedOrder)

	//not found
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.String(http.StatusOK, "Under construct")
	})

	// middlewares
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ownMiddleware.Auth(sm))

	return e
}
