package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/marketplace/orders"
	"mephiMainProject/pkg/services/marketplace/product"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/profile"
	"mephiMainProject/pkg/services/server/session"
	"mephiMainProject/pkg/services/server/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MarketplaceHandler struct {
	Logger             *zap.SugaredLogger
	CurrentCfg         *config.Config
	MarketPlaceManager product.MarketplaceServiceClient
	OrdersManager      orders.OrderServiceClient
	ProfileRepo        profile.ProfileRepo
}

func (mh *MarketplaceHandler) GetProduct(c echo.Context) error {
	formData := NewFormData()

	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page", formData)
	}
	formData.Values["username"] = currentSession.Username

	productId := c.Param("id")
	currentProduct, err := mh.MarketPlaceManager.GetProduct(context.Background(), &product.ProductID{ProductID: productId})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page", formData)
	}
	if currentProduct.OwnerUsername == currentSession.Username {
		formData.Values["currentUserIsOwner"] = "1"
	} else {
		_, err = mh.ProfileRepo.GetProfile(currentSession.Username)
		if errors.Is(err, nil) {
			formData.Values["currentUserIsNotOwner"] = "1"
		}
	}
	formData.Values["id"] = productId
	formData.Values["name"] = currentProduct.Name
	formData.Values["description"] = currentProduct.Description
	formData.Values["price"] = strconv.Itoa(int(currentProduct.Price))
	formData.Values["mainPhoto"] = currentProduct.MainPhoto
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
	currentSession, err := session.SessionFromContext(c)
	formData := NewFormData()
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page-data", formData)
	}
	productId := c.Param("id")
	productFromDB, err := mh.MarketPlaceManager.GetProduct(context.Background(), &product.ProductID{ProductID: productId})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page-data", formData)
	}
	if productFromDB.OwnerUsername == currentSession.Username {
		_, err = mh.MarketPlaceManager.DeleteProduct(context.Background(), &product.ProductID{ProductID: productId})
		if err != nil {
			formData.Errors["error"] = err.Error()
			return c.Render(422, "marketplace-item-page-data", formData)
		}
		err = utils.DeleteFile(productFromDB.PhotoUrls)
		if err != nil {
			formData.Errors["error"] = err.Error()
			return c.Render(422, "marketplace-item-page-data", formData)
		}
	} else {
		formData.Errors["error"] = "You aren't owner of this product"
		return c.Render(422, "marketplace-item-page-data", formData)
	}
	return c.String(http.StatusOK, "Successfully deleted product")
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
	}
	price, err := strconv.Atoi(c.FormValue("price"))
	newProduct.Price = int64(price)

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
	newProduct.MainPhoto = photoUrls[0]

	_, err = mh.MarketPlaceManager.CreateProduct(context.Background(), newProduct)
	if err != nil {
		fmt.Printf("CreateProduct err - %v", err)
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-form-add", formData)
	}
	return c.Redirect(http.StatusSeeOther, "/marketplace/products/"+currentSession.Username)
}

// Orders

func (mh *MarketplaceHandler) GetOrders(c echo.Context) error {
	formData := NewFormData()
	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "orders-view", formData)
	}
	formData.Values["username"] = currentSession.Username
	allOrders, err := mh.OrdersManager.GetUserOrders(context.Background(), &orders.Buyer{BuyerUsername: currentSession.Username})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "orders-view", formData)
	}

	if len(allOrders.GetOrders()) == 0 {
		formData.Values["empty"] = "empty"
	} else {
		ordersToReturn := allOrders.GetOrders()
		for i, v := range ordersToReturn {
			productName, err := mh.MarketPlaceManager.GetProduct(context.Background(), &product.ProductID{ProductID: strconv.Itoa(int(v.ProductId))})
			if err != nil {
				continue
			}
			ordersToReturn[i].BuyerUsername = productName.Name
		}
		formData.Orders = ordersToReturn
	}
	return c.Render(http.StatusOK, "orders-view", formData)
}

func (mh *MarketplaceHandler) GetSales(c echo.Context) error {
	formData := NewFormData()
	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "orders-view", formData)
	}
	formData.Values["username"] = currentSession.Username
	allSales, err := mh.OrdersManager.GetSellerOrders(context.Background(), &orders.Seller{SellerUsername: currentSession.Username})
	if err != nil {
		formData.Values["error"] = err.Error()
		mh.Logger.Infof("GetSeller orders handler err - %v\n", err.Error())
		return c.Render(422, "sales-view", formData)
	}
	if len(allSales.GetOrders()) == 0 {
		formData.Values["empty"] = "empty"
	} else {
		ordersToReturn := allSales.GetOrders()
		for i, v := range ordersToReturn {
			productName, err := mh.MarketPlaceManager.GetProduct(context.Background(), &product.ProductID{ProductID: strconv.Itoa(int(v.ProductId))})
			if err != nil {
				continue
			}
			ordersToReturn[i].SellerUsername = productName.Name
		}
		formData.Orders = ordersToReturn
	}
	return c.Render(http.StatusOK, "sales-view", formData)
}

func (mh *MarketplaceHandler) ProceedOrder(c echo.Context) error {
	formData := NewFormData()
	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page", formData)
	}
	formData.Values["username"] = currentSession.Username
	currentUser, err := mh.ProfileRepo.GetProfile(currentSession.Username)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "marketplace-item-page", formData)
	}
	buyerName := c.FormValue("buyerName")
	if buyerName == "" {
		formData.Errors["error"] = "Некорректные данные (поле 'Имя') при оформлении заказе"
		return c.Render(422, "marketplace-item-page", formData)
	}
	buyerAddress := c.FormValue("address")
	if buyerAddress == "" {
		buyerAddress = currentUser.Address
	}
	contacts := c.FormValue("contacts")
	if contacts == "" {
		mh.Logger.Infof("Incorrect input by contacts\n")
		formData.Errors["error"] = "Некорректные данные (поле 'Контакты') при оформлении заказе"
		return c.Render(422, "marketplace-item-page", formData)
	}
	desc := c.FormValue("description")
	if desc == "" {
		desc = "Пусто" + " | Указанный контакт для связи: " + contacts
	} else {
		desc = desc + " | Указанный контакт для связи: " + contacts
	}
	productId := c.FormValue("ProductID")
	intProductId, err := strconv.ParseUint(productId, 10, 64)
	if err != nil {
		formData.Errors["error"] = "Error by parsing productID from headers"
		return c.Render(422, "marketplace-item-page", formData)
	}
	currentProduct, err := mh.MarketPlaceManager.GetProduct(context.Background(), &product.ProductID{ProductID: productId})
	if err != nil {
		mh.Logger.Infof("Can't get product with this ID\n")
		formData.Errors["error"] = fmt.Sprintf("Error by getting product with ID %v from DB\n", productId)
		return c.Render(422, "marketplace-item-page", formData)
	}

	if currentProduct.OwnerUsername == currentSession.Username {
		formData.Values["currentUserIsOwner"] = "1"
	}
	formData.Values["id"] = productId
	formData.Values["name"] = currentProduct.Name
	formData.Values["description"] = currentProduct.Description
	formData.Values["price"] = strconv.Itoa(int(currentProduct.Price))
	formData.Values["mainPhoto"] = currentProduct.MainPhoto
	formData.Values["photoUrls"] = strings.Join(currentProduct.PhotoUrls, " | ")

	_, err = mh.OrdersManager.CreateOrder(context.Background(), &orders.Order{
		SellerUsername: currentProduct.OwnerUsername,
		BuyerUsername:  currentSession.Username,
		BuyerName:      buyerName,
		ProductId:      int64(intProductId),
		ProductCount:   1,
		OrderComment:   desc,
		OrderAddress:   buyerAddress,
	})

	if err != nil {
		mh.Logger.Infof("Error by creating new order: %v\n", err.Error())
		formData.Errors["error"] = fmt.Sprintf("Error by creating new order: %v\n", err.Error())
		return c.Render(422, "marketplace-item-page", formData)
	}
	return c.Redirect(http.StatusSeeOther, "/marketplace/orders/")
}

func (mh *MarketplaceHandler) AcceptOrder(c echo.Context) error {
	formData := NewFormData()
	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "sales-view", formData)
	}
	formData.Values["username"] = currentSession.Username
	currentOrderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	currentOrder, err := mh.OrdersManager.GetOrder(context.Background(), &orders.OrderID{Id: currentOrderID})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "sales-view", formData)
	}
	if currentOrder.SellerUsername != currentSession.Username {
		formData.Errors["error"] = "You are not the seller of this order\n"
		return c.Render(422, "sales-view", formData)
	}

	_, err = mh.OrdersManager.AcceptOrder(context.Background(), &orders.OrderID{Id: currentOrderID})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(500, "sales-view", formData)
	}
	return c.Render(http.StatusOK, "sale", formData)
}

func (mh *MarketplaceHandler) CompleteOrder(c echo.Context) error {
	formData := NewFormData()
	currentSession, err := session.SessionFromContext(c)
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "sales-view", formData)
	}
	formData.Values["username"] = currentSession.Username
	currentOrderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	currentOrder, err := mh.OrdersManager.GetOrder(context.Background(), &orders.OrderID{Id: currentOrderID})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(422, "sales-view", formData)
	}
	if currentOrder.SellerUsername != currentSession.Username {
		formData.Errors["error"] = "You are not the seller of this order\n"
		return c.Render(422, "sales-view", formData)
	}

	_, err = mh.OrdersManager.CompleteOrder(context.Background(), &orders.OrderID{Id: currentOrderID})
	if err != nil {
		formData.Errors["error"] = err.Error()
		return c.Render(500, "sales-view", formData)
	}
	return c.Render(http.StatusOK, "sale", formData)
}
