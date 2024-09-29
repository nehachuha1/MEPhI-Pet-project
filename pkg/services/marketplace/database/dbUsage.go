package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"mephiMainProject/pkg/services/marketplace/config"
	"strings"
)

var (
	c                           redis.Conn
	ordersKey                   = "ORDERSCREATED"
	errorIncrementOrdersCounter = errors.New("error while incrementing orders counter")
)

type DatabaseORM struct {
	Pgx *config.PostgreDB
	Rds *config.RedisDB
}

var _ DatabaseControl = &DatabaseORM{} // check DatabaseControl interface implementation

func NewPgxConn(cfg *config.Config) *config.PostgreDB {
	dsn := "postgres://" + cfg.Database.PgxUser + ":" + cfg.Database.PgxPassword + "@"
	dsn = dsn + cfg.Database.PgxAddress + ":" + cfg.Database.PgxPort + "/"
	dsn = dsn + cfg.Database.PgxDB

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Postgre connection err - %v\n", err)
		return nil
	}
	return &config.PostgreDB{
		DB: db,
	}
}

func NewDBUsage(cfg *config.Config) *DatabaseORM {
	return &DatabaseORM{
		Pgx: NewPgxConn(cfg),
		Rds: &config.RedisDB{RedisConnection: redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(cfg.Database.RedisURL)
			},
			MaxIdle:     8,
			MaxActive:   0,
			IdleTimeout: 100,
		}},
	}
}

func (db *DatabaseORM) GetAllProducts() ([]config.Product, error) {
	rows, err := db.Pgx.DB.Query("SELECT id, name, owner_username, price, description, create_date, edit_date, is_active, views, photo_urls, main_photo FROM public.products;")
	if err != nil {
		log.Printf("GetAllProducts err - %v\n", err)
		return []config.Product{}, err
	}

	var allProducts []config.Product
	var photoUrls string
	for rows.Next() {
		var currentProduct config.Product
		err = rows.Scan(&currentProduct.ID, &currentProduct.Name, &currentProduct.OwnerUsername, &currentProduct.Price, &currentProduct.Description, &currentProduct.CreateDate,
			&currentProduct.EditDate, &currentProduct.IsActive, &currentProduct.Views, &photoUrls, &currentProduct.MainPhoto)
		if err != nil {
			log.Printf("Error while scanning current product - %v\n", err)
		}
		normalizedPhotoURLs := strings.Split(photoUrls[1:len(photoUrls)-1], ",")
		currentProduct.PhotoURLs = normalizedPhotoURLs
		allProducts = append(allProducts, currentProduct)
	}
	return allProducts, nil
}

func (db *DatabaseORM) GetProduct(productID string) (config.Product, error) {
	rows, err := db.Pgx.DB.Query("SELECT name, owner_username, price, description, create_date, edit_date, is_active, views, photo_urls, main_photo FROM public.products WHERE id=$1;",
		productID)
	if err != nil {
		log.Printf("GetAllProducts err - %v\n", err)
		return config.Product{}, err
	}

	var currentProduct config.Product
	var photoUrls string
	for rows.Next() {
		err = rows.Scan(&currentProduct.Name, &currentProduct.OwnerUsername, &currentProduct.Price, &currentProduct.Description, &currentProduct.CreateDate,
			&currentProduct.EditDate, &currentProduct.IsActive, &currentProduct.Views, &photoUrls, &currentProduct.MainPhoto)
		if err != nil {
			log.Printf("Error while scanning current product - %v", err)
			return config.Product{}, err
		}
	}
	normalizedPhotoURLs := strings.Split(photoUrls[1:len(photoUrls)-1], ",")
	currentProduct.PhotoURLs = normalizedPhotoURLs
	return currentProduct, nil
}

func (db *DatabaseORM) CreateProduct(product config.Product) error {
	_, err := db.Pgx.DB.Exec("INSERT INTO public.products(name, owner_username, price, description, create_date, edit_date, is_active, views, photo_urls, main_photo) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		&product.Name, &product.OwnerUsername, &product.Price, &product.Description, &product.CreateDate,
		&product.EditDate, &product.IsActive, &product.Views, &product.PhotoURLs, &product.MainPhoto,
	)

	if err != nil {
		log.Printf("Create product err - %v", err)
		return err
	}

	return nil
}

func (db *DatabaseORM) EditProduct(product config.Product, productID string) error {
	_, err := db.Pgx.DB.Exec("UPDATE public.products SET name=$2, owner_username=$3, price=$4, description=$5, create_date=$6, edit_date=$7, is_active=$8, views=$9, photo_urls=$10 WHERE id=$1;",
		productID, &product.Name, &product.OwnerUsername, &product.Price, &product.Description, &product.CreateDate,
		&product.EditDate, &product.IsActive, &product.Views, &product.PhotoURLs,
	)
	if err != nil {
		log.Printf("Edit product err - %v", err)
		return err
	}

	return nil
}

func (db *DatabaseORM) DeleteProduct(productID string) error {
	_, err := db.Pgx.DB.Exec("DELETE FROM public.products WHERE id=$1;", productID)

	if err != nil {
		log.Printf("Delete product err - %v", err)
		return err
	}

	return nil
}

// Orders service

func (db *DatabaseORM) GetSellerOrders(seller *config.Seller) (*config.AllOrders, error) {
	rows, err := db.Pgx.DB.Query(getSellerOrdersScript, &seller.SellerUsername)
	if err != nil {
		return &config.AllOrders{}, err
	}

	sellerOrders := &config.AllOrders{}
	for rows.Next() {
		currentOrder := &config.Order{}
		err = rows.Scan(&currentOrder.Id, &currentOrder.SellerUsername, &currentOrder.BuyerUsername,
			&currentOrder.ProductId, &currentOrder.ProductCount, &currentOrder.OrderComment, &currentOrder.OrderStatus, &currentOrder.IsCompleted)
		if err != nil {
			return &config.AllOrders{}, err
		}
		sellerOrders.Orders = append(sellerOrders.Orders, currentOrder)
	}
	return sellerOrders, nil
}

func (db *DatabaseORM) GetUserOrders(buyer *config.Buyer) (*config.AllOrders, error) {
	rows, err := db.Pgx.DB.Query(getUserOrdersScript, &buyer.BuyerUsername)
	if err != nil {
		log.Printf("DB Error 1: %s\n", err.Error())
		return &config.AllOrders{}, err
	}

	userOrders := &config.AllOrders{}
	for rows.Next() {
		currentOrder := &config.Order{}
		err = rows.Scan(&currentOrder.Id, &currentOrder.SellerUsername, &currentOrder.BuyerUsername,
			&currentOrder.ProductId, &currentOrder.ProductCount, &currentOrder.OrderComment, &currentOrder.OrderStatus, &currentOrder.IsCompleted)
		if err != nil {
			log.Printf("DB Error 2: %s\n", err.Error())
			//return &config.AllOrders{}, err
		}
		userOrders.Orders = append(userOrders.Orders, currentOrder)
	}
	return userOrders, nil
}
func (db *DatabaseORM) CreateOrder(order *config.Order) (*config.OrderID, error) {
	_, err := db.Pgx.DB.ExecContext(context.Background(), createOrderScript, order.SellerUsername, order.BuyerUsername, order.ProductId,
		order.ProductCount, order.OrderComment, order.OrderAddress)
	if err != nil {
		return &config.OrderID{}, err
	}

	lastOrderNum, err := redis.Int(c.Do("GET", ordersKey))
	if errors.Is(err, redis.ErrNil) {
		lastOrderNum, err = redis.Int(c.Do("SET", ordersKey, 1))
		if err != nil {
			return &config.OrderID{}, err
		}
	}
	_, err = redis.Int(c.Do("INCRBY", ordersKey, 1))
	if err != nil {
		return &config.OrderID{Id: int64(lastOrderNum)}, errorIncrementOrdersCounter
	}
	return &config.OrderID{Id: int64(lastOrderNum)}, nil
}

func (db *DatabaseORM) GetOrder(orderId *config.OrderID) (*config.Order, error) {
	rows, err := db.Pgx.DB.QueryContext(context.Background(), getOrderScript, orderId.Id)
	if err != nil {
		return &config.Order{}, err
	}

	currentOrder := &config.Order{}
	for rows.Next() {
		err = rows.Scan(&currentOrder.Id, &currentOrder.SellerUsername, &currentOrder.BuyerUsername, &currentOrder.ProductId,
			&currentOrder.ProductCount, &currentOrder.OrderComment, &currentOrder.OrderAddress, &currentOrder.OrderAddress,
			&currentOrder.IsCompleted)
		if err != nil {
			return &config.Order{}, err
		}
	}
	return currentOrder, nil
}

func (db *DatabaseORM) AcceptOrder(orderId *config.OrderID) (*config.Response, error) {
	_, err := db.Pgx.DB.ExecContext(context.Background(), acceptOrderScript, orderId.Id)

	if err != nil {
		return &config.Response{}, err
	}

	resp := &config.Response{
		Code:    200,
		Message: "Successfully accepted order",
	}
	return resp, nil
}

func (db *DatabaseORM) CompleteOrder(orderId *config.OrderID) (*config.Response, error) {
	_, err := db.Pgx.DB.ExecContext(context.Background(), completeOrderScript, orderId.Id)

	if err != nil {
		return &config.Response{}, err
	}

	resp := &config.Response{
		Code:    200,
		Message: "Successfully completed order",
	}
	return resp, nil
}

func (db *DatabaseORM) CheckUserBlock(user *config.User) (*config.UserBlock, error) {
	rows, err := db.Pgx.DB.QueryContext(context.Background(), checkUserBlockScript, user.Username)
	if err != nil {
		return &config.UserBlock{}, err
	}

	block := &config.UserBlock{}
	for rows.Next() {
		err = rows.Scan(&block.Id, &block.IntruderUsername, &block.ModeratorUsername,
			&block.BanReason, &block.BanDate, &block.ExpiresAt)
		if err != nil {
			return &config.UserBlock{}, err
		}
	}
	return block, nil
}

func (db *DatabaseORM) BlockUser(blockInfo *config.UserBlock) (*config.Response, error) {
	_, err := db.Pgx.DB.ExecContext(context.Background(), blockUserScript, blockInfo.IntruderUsername,
		blockInfo.ModeratorUsername, blockInfo.BanReason, blockInfo.BanDate, blockInfo.ExpiresAt)

	if err != nil {
		return &config.Response{}, err
	}

	resp := &config.Response{
		Code: 200,
		Message: fmt.Sprintf("Blocked username %s | Banned by: %s | Reason: %s | Expires at %s\n", blockInfo.IntruderUsername,
			blockInfo.ModeratorUsername, blockInfo.BanReason, blockInfo.ExpiresAt),
	}
	return resp, nil
}

func (db *DatabaseORM) UnblockUser(user *config.User) (*config.Response, error) {
	_, err := db.Pgx.DB.ExecContext(context.Background(), unblockUserScript, user.Username)

	if err != nil {
		return &config.Response{}, err
	}

	resp := &config.Response{
		Code:    200,
		Message: fmt.Sprintf("Unblocked username %s\n", user.Username),
	}
	return resp, nil
}
