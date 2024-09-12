package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/marketplace/product"
	"mephiMainProject/pkg/services/server/session"
	"net/http"
	"strconv"
	"time"
)

type MarketplaceHandler struct {
	Logger             *zap.SugaredLogger
	MarketPlaceManager product.MarketplaceServiceClient
}

func (mh *MarketplaceHandler) GetProducts(c echo.Context) error {
	allProducts, err := mh.MarketPlaceManager.GetAllProducts(context.Background(), &product.Nothing{})
	currentSession, _ := session.SessionFromContext(c)

	formData := NewFormData()
	formData.Values["username"] = currentSession.Username
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(http.StatusOK, "marketplace-view", formData)
	}
	formData.Products = allProducts.GetProducts()
	formData.Values["marketplace"] = "marketplace"
	return c.Render(http.StatusOK, "marketplace-view", formData)
}

func (mh *MarketplaceHandler) CreateProductGet(c echo.Context) error {
	formData := NewFormData()
	return c.Render(200, "marketplace-form-add", formData)
}

func (mh *MarketplaceHandler) CreateProductPost(c echo.Context) error {
	formData := NewFormData()
	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-form-add", formData)
	}
	newProduct := &product.Product{
		Name:          c.FormValue("name"),
		OwnerUsername: currentSession.Username,
		Description:   c.FormValue("description"),
		CreateDate:    time.Now().Format("01-02-2006 15:04:05"),
		EditDate:      time.Now().Format("01-02-2006 15:04:05"),
		IsActive:      true,
		Views:         1,
		PhotoUrls:     []string{"123", "233", "556"},
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	newProduct.Price = int64(price)
	_, err = mh.MarketPlaceManager.CreateProduct(context.Background(), newProduct)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-form-add", formData)
	}
	return c.Redirect(http.StatusSeeOther, "/marketplace/products/"+currentSession.Username)
}
