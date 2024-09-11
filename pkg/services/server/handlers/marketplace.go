package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/marketplace/product"
	"net/http"
)

type MarketplaceHandler struct {
	Logger             *zap.SugaredLogger
	MarketPlaceManager product.MarketplaceServiceClient
}

func (mh *MarketplaceHandler) GetProducts(c echo.Context) error {
	allProducts, err := mh.MarketPlaceManager.GetAllProducts(context.Background(), &product.Nothing{})

	formData := NewFormData()
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(http.StatusOK, "marketplace-table", formData)
	}
	formData.Products = allProducts.GetProducts()
	formData.Values["marketplace"] = "marketplace"
	return c.Render(http.StatusOK, "marketplace-table", formData)
}
