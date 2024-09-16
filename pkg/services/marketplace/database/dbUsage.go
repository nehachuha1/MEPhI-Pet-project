package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"mephiMainProject/pkg/services/marketplace/config"
	"strings"
)

type DatabaseORM struct {
	Pgx *config.PostgreDB
}

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
	}
}

func (db *DatabaseORM) GetAllProducts() ([]config.Product, error) {
	rows, err := db.Pgx.DB.Query("SELECT id, name, owner_username, price, description, create_date, edit_date, is_active, views FROM public.products;")
	if err != nil {
		log.Printf("GetAllProducts err - %v\n", err)
		return []config.Product{}, err
	}

	var allProducts []config.Product
	for rows.Next() {
		var currentProduct config.Product
		err = rows.Scan(&currentProduct.ID, &currentProduct.Name, &currentProduct.OwnerUsername, &currentProduct.Price, &currentProduct.Description, &currentProduct.CreateDate,
			&currentProduct.EditDate, &currentProduct.IsActive, &currentProduct.Views)
		if err != nil {
			log.Printf("Error while scanning current product - %v\n", err)
		}

		allProducts = append(allProducts, currentProduct)
	}
	return allProducts, nil
}

func (db *DatabaseORM) GetProduct(productID string) (config.Product, error) {
	rows, err := db.Pgx.DB.Query("SELECT name, owner_username, price, description, create_date, edit_date, is_active, views, photo_urls FROM public.products WHERE id=$1;",
		productID)
	if err != nil {
		log.Printf("GetAllProducts err - %v\n", err)
		return config.Product{}, err
	}

	var currentProduct config.Product
	var photoUrls string
	for rows.Next() {
		err = rows.Scan(&currentProduct.Name, &currentProduct.OwnerUsername, &currentProduct.Price, &currentProduct.Description, &currentProduct.CreateDate,
			&currentProduct.EditDate, &currentProduct.IsActive, &currentProduct.Views, &photoUrls)
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
	_, err := db.Pgx.DB.Exec("INSERT INTO public.products(name, owner_username, price, description, create_date, edit_date, is_active, views, photo_urls) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
		&product.Name, &product.OwnerUsername, &product.Price, &product.Description, &product.CreateDate,
		&product.EditDate, &product.IsActive, &product.Views, &product.PhotoURLs,
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
