package config

import "database/sql"

type Product struct {
	ID            int64
	Name          string
	OwnerUsername string
	Price         int64
	Description   string
	CreateDate    string
	EditDate      string
	IsActive      bool
	Views         int64
	PhotoURLs     []string
	MainPhoto     string
}

type PostgreDB struct {
	DB *sql.DB
}
