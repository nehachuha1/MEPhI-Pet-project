package handlers

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/marketplace/product"
	"mephiMainProject/pkg/services/server/session"
	"mephiMainProject/pkg/services/server/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MarketplaceHandler struct {
	Logger             *zap.SugaredLogger
	MarketPlaceManager product.MarketplaceServiceClient
}

func (mh *MarketplaceHandler) GetProduct(c echo.Context) error {
	formData := NewFormData()

	currentUser, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page", formData)
	}
	formData.Values["username"] = currentUser.Username

	productId := c.Param("id")
	currentProduct, err := mh.MarketPlaceManager.GetProduct(context.Background(), &product.ProductID{ProductID: productId})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page", formData)
	}
	if currentProduct.OwnerUsername == currentUser.Username {
		formData.Values["currentUserIsOwner"] = "1"
	}
	formData.Values["id"] = productId
	formData.Values["name"] = currentProduct.Name
	formData.Values["description"] = currentProduct.Description
	formData.Values["price"] = strconv.Itoa(int(currentProduct.Price))
	formData.Values["photoUrls"] = strings.Join(currentProduct.PhotoUrls, " | ")

	return c.Render(http.StatusOK, "marketplace-item-page", formData)
}

func (mh *MarketplaceHandler) GetProducts(c echo.Context) error {
	formData := NewFormData()
	allProducts, err := mh.MarketPlaceManager.GetAllProducts(context.Background(), &product.Nothing{})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-view", formData)
	}

	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(http.StatusOK, "marketplace-view", formData)
	}

	formData.Values["username"] = currentSession.Username
	formData.Products = allProducts.GetProducts()
	formData.Values["marketplace"] = "marketplace"
	return c.Render(http.StatusOK, "marketplace-view", formData)
}

func (mh *MarketplaceHandler) DeleteProduct(c echo.Context) error {
	return c.String(http.StatusOK, "Delete product realization...")
}

func (mh *MarketplaceHandler) GetUserProducts(c echo.Context) error {
	allProducts, err := mh.MarketPlaceManager.GetAllProducts(context.Background(), &product.Nothing{})
	currentSession, _ := session.SessionFromContext(c)
	username := c.Param("username")

	formData := NewFormData()
	if currentSession.Username == username {
		formData.Values["currentUserIsOwner"] = currentSession.Username
	}
	formData.Values["username"] = currentSession.Username
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(http.StatusOK, "marketplace-view", formData)
	}

	var returnProducts []*product.Product
	products := allProducts.GetProducts()
	for _, pr := range products {
		if pr.OwnerUsername == username {
			returnProducts = append(returnProducts, pr)
		}
	}
	formData.Products = returnProducts
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

	// обработка фотографий при загрузке
	form, err := c.MultipartForm()
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-form-add", formData)
	}
	files := form.File["files"]
	photoUrls, err := utils.ServeFiles(files)
	if err != nil {
		fmt.Printf("ServeFiels err - %v", err)
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-form-add", formData)
	}
	newProduct.PhotoUrls = photoUrls

	_, err = mh.MarketPlaceManager.CreateProduct(context.Background(), newProduct)
	if err != nil {
		fmt.Printf("CreateProduct err - %v", err)
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-form-add", formData)
	}
	return c.Redirect(http.StatusSeeOther, "/marketplace/products/"+currentSession.Username)
}
