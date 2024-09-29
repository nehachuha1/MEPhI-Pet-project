package database

import (
	"mephiMainProject/pkg/services/marketplace/config"
)

type DatabaseControl interface {
	// Products functions

	GetAllProducts() ([]config.Product, error)
	GetProduct(productID string) (config.Product, error)
	CreateProduct(product config.Product) error
	EditProduct(product config.Product, productID string) error
	DeleteProduct(productID string) error

	// Orders functions

	GetSellerOrders(seller *config.Seller) (*config.AllOrders, error)
	GetUserOrders(buyer *config.Buyer) (*config.AllOrders, error)
	CreateOrder(order *config.Order) (*config.OrderID, error)
	GetOrder(orderId *config.OrderID) (*config.Order, error)
	AcceptOrder(orderId *config.OrderID) (*config.Response, error)
	CompleteOrder(orderId *config.OrderID) (*config.Response, error)
	CheckUserBlock(user *config.User) (*config.UserBlock, error)
	BlockUser(blockInfo *config.UserBlock) (*config.Response, error)
	UnblockUser(user *config.User) (*config.Response, error)
}
