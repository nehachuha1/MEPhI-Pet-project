package config

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
)

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

type Seller struct {
	Id                int64
	SellerUsername    string
	Accepted          bool
	ModeratorUsername string
	IsActive          bool
	IsBanned          bool
	BanId             int64
	Balance           int64
	Transactions      []string
}

type Buyer struct {
	Id            int64
	BuyerUsername string
}

type User struct {
	Username string
}

type UserBlock struct {
	Id                int64
	IntruderUsername  string
	ModeratorUsername string
	BanReason         string
	BanDate           string
	ExpiresAt         string
}

type Response struct {
	Code    int64
	Message string
}

type Order struct {
	Id             int64
	SellerUsername string
	BuyerUsername  string
	BuyerName      string
	ProductId      int64
	ProductCount   int64
	OrderComment   string
	OrderAddress   string
	OrderStatus    string
	IsCompleted    bool
}

type OrderID struct {
	Id int64
}

type AllOrders struct {
	Orders []*Order `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
	Page   int64    `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
}

type PostgreDB struct {
	DB *sql.DB
}

type RedisDB struct {
	RedisConnection redis.Pool
}
