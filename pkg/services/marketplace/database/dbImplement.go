package database

import (
	"mephiMainProject/pkg/services/marketplace/config"
)

type DatabaseControl interface {
	GetAllProducts() ([]config.Product, error)
	GetProduct(productID string) (config.Product, error)
	CreateProduct(product config.Product) error
	EditProduct(product config.Product, productID string) error
	DeleteProduct(productID string) error
}
