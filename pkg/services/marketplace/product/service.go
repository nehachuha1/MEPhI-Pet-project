package product

import (
	"context"
	"log"
	"mephiMainProject/pkg/services/marketplace/config"
	"mephiMainProject/pkg/services/marketplace/database"
	"strconv"
)

type MarketplaceService struct {
	UnimplementedMarketplaceServiceServer

	Database *database.DatabaseORM
}

func NewMarketplaceService(cfg *config.Config) *MarketplaceService {
	return &MarketplaceService{
		Database: database.NewDBUsage(cfg),
	}
}

func (ms *MarketplaceService) GetProduct(ctx context.Context, pID *ProductID) (*Product, error) {
	currentProduct, err := ms.Database.GetProduct(pID.ProductID)
	if err != nil {
		log.Printf("Get product err in marketplace module - %v", err)
		return &Product{}, err
	}

	returnProduct := &Product{
		Id:            currentProduct.ID,
		Name:          currentProduct.Name,
		OwnerUsername: currentProduct.OwnerUsername,
		Price:         currentProduct.Price,
		Description:   currentProduct.Description,
		CreateDate:    currentProduct.CreateDate,
		EditDate:      currentProduct.EditDate,
		IsActive:      currentProduct.IsActive,
		Views:         currentProduct.Views,
		PhotoUrls:     currentProduct.PhotoURLs,
	}

	return returnProduct, nil
}

func (ms *MarketplaceService) GetAllProducts(ctx context.Context, nth *Nothing) (*AllProducts, error) {
	allProducts, err := ms.Database.GetAllProducts()
	if err != nil {
		return &AllProducts{}, err
	}

	var returnProducts AllProducts
	for _, elem := range allProducts {
		var currentProduct *Product

		currentProduct = &Product{
			Id:            elem.ID,
			Name:          elem.Name,
			OwnerUsername: elem.OwnerUsername,
			Price:         elem.Price,
			Description:   elem.Description,
			CreateDate:    elem.CreateDate,
			EditDate:      elem.EditDate,
			IsActive:      elem.IsActive,
			Views:         elem.Views,
			PhotoUrls:     elem.PhotoURLs,
		}
		returnProducts.Products = append(returnProducts.Products, currentProduct)
	}

	return &returnProducts, nil
}

func (ms *MarketplaceService) CreateProduct(ctx context.Context, p *Product) (*Response, error) {
	tempVar := config.Product{
		Name:          p.Name,
		OwnerUsername: p.OwnerUsername,
		Price:         p.Price,
		Description:   p.Description,
		CreateDate:    p.CreateDate,
		EditDate:      p.EditDate,
		IsActive:      p.IsActive,
		Views:         p.Views,
		PhotoURLs:     p.GetPhotoUrls(),
	}
	err := ms.Database.CreateProduct(tempVar)
	if err != nil {
		log.Printf("Create product err in marketplace module - %v", err)
		return &Response{
			Code:    500,
			Message: err.Error(),
		}, err
	}

	return &Response{Code: 200, Message: "OK"}, nil
}

func (ms *MarketplaceService) EditProduct(ctx context.Context, p *Product) (*Response, error) {
	currentProduct := config.Product{
		ID:            p.Id,
		Name:          p.Name,
		OwnerUsername: p.OwnerUsername,
		Price:         p.Price,
		Description:   p.Description,
		CreateDate:    p.CreateDate,
		EditDate:      p.EditDate,
		IsActive:      p.IsActive,
		Views:         p.Views,
		PhotoURLs:     p.PhotoUrls,
	}
	err := ms.Database.EditProduct(currentProduct, strconv.FormatInt(p.Id, 10))

	if err != nil {
		log.Printf("Edit product err in marketplace module - %v", err)
		return &Response{
			Code:    500,
			Message: err.Error(),
		}, err
	}
	return &Response{Code: 200, Message: "OK"}, nil
}

func (ms *MarketplaceService) DeleteProduct(ctx context.Context, pID *ProductID) (*Response, error) {
	err := ms.Database.DeleteProduct(pID.ProductID)

	if err != nil {
		log.Printf("Delete product err in marketplace module - %v", err)
		return &Response{
			Code:    500,
			Message: err.Error(),
		}, err
	}
	return &Response{Code: 200, Message: "OK"}, nil
}
