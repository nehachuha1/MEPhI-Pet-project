package orders

import (
	"context"
	"errors"
	"mephiMainProject/pkg/services/marketplace/config"
	"mephiMainProject/pkg/services/marketplace/database"
)

var (
	emptySellerUsername  = errors.New("empty seller username")
	emptyBuyerUsername   = errors.New("empty buyer username")
	wrongOrderStruct     = errors.New("current Order struct is empty or some of fields may have zero values")
	wrongUserStruct      = errors.New("current User struct is empty or some of fields may have zero values")
	wrongUserBlockStruct = errors.New("current UserBlock struct is empty or some of fields may have zero values")
)

type OrderService struct {
	UnimplementedOrderServiceServer

	Database *database.DatabaseORM
}

func NewOrderService(cfg *config.Config) *OrderService {
	return &OrderService{
		Database: database.NewDBUsage(cfg),
	}
}

func (ors *OrderService) GetSellerOrders(ctx context.Context, seller *Seller) (*AllOrders, error) {
	if seller.SellerUsername == "" {
		return &AllOrders{}, emptySellerUsername
	}
	currentSeller := &config.Seller{
		Id:                seller.Id,
		SellerUsername:    seller.SellerUsername,
		Accepted:          seller.Accepted,
		ModeratorUsername: seller.ModeratorUsername,
		IsActive:          seller.IsActive,
		IsBanned:          seller.IsBanned,
		BanId:             seller.BanId,
		Balance:           seller.Balance,
		Transactions:      seller.Transactions,
	}
	res, err := ors.Database.GetSellerOrders(currentSeller)
	if err != nil {
		return &AllOrders{}, err
	}
	orders := make([]*Order, 0)
	for _, ord := range res.Orders {
		currentOrder := &Order{
			Id:             ord.Id,
			SellerUsername: ord.SellerUsername,
			BuyerUsername:  ord.BuyerUsername,
			BuyerName:      ord.BuyerName,
			ProductId:      ord.ProductId,
			ProductCount:   ord.ProductCount,
			OrderComment:   ord.OrderComment,
			OrderAddress:   ord.OrderAddress,
			OrderStatus:    ord.OrderStatus,
			IsCompleted:    ord.IsCompleted,
		}
		orders = append(orders, currentOrder)
	}
	return &AllOrders{Orders: orders, Page: res.Page}, nil
}

func (ors *OrderService) GetUserOrders(ctx context.Context, buyer *Buyer) (*AllOrders, error) {
	if !validateBuyer(buyer) {
		return &AllOrders{}, emptyBuyerUsername
	}

	res, err := ors.Database.GetUserOrders(&config.Buyer{BuyerUsername: buyer.BuyerUsername})
	if err != nil {
		return &AllOrders{}, err
	}

	returnOrders := &AllOrders{Page: 1}
	for _, ord := range res.Orders {
		currentOrder := &Order{
			Id:             ord.Id,
			SellerUsername: ord.SellerUsername,
			BuyerUsername:  ord.BuyerUsername,
			BuyerName:      ord.BuyerName,
			ProductId:      ord.ProductId,
			ProductCount:   ord.ProductCount,
			OrderComment:   ord.OrderComment,
			OrderAddress:   ord.OrderAddress,
			OrderStatus:    ord.OrderStatus,
			IsCompleted:    ord.IsCompleted,
		}
		returnOrders.Orders = append(returnOrders.Orders, currentOrder)
	}
	return returnOrders, nil
}

func (ors *OrderService) CreateOrder(ctx context.Context, order *Order) (*OrderID, error) {
	if !validateOrder(order) {
		return &OrderID{}, wrongOrderStruct
	}
	currentOrder := &config.Order{
		SellerUsername: order.SellerUsername,
		BuyerUsername:  order.BuyerUsername,
		BuyerName:      order.BuyerName,
		ProductId:      order.ProductId,
		ProductCount:   order.ProductCount,
		OrderComment:   order.OrderComment,
		OrderAddress:   order.OrderAddress,
	}
	res, err := ors.Database.CreateOrder(currentOrder)
	if err != nil {
		return &OrderID{}, err
	}
	return &OrderID{Id: res.Id}, nil
}

func (ors *OrderService) GetOrder(ctx context.Context, orderID *OrderID) (*Order, error) {
	res, err := ors.Database.GetOrder(&config.OrderID{Id: orderID.Id})
	if err != nil {
		return &Order{}, err
	}
	currentOrder := &Order{
		Id:             res.Id,
		SellerUsername: res.SellerUsername,
		BuyerUsername:  res.BuyerUsername,
		BuyerName:      res.BuyerName,
		ProductId:      res.ProductId,
		ProductCount:   res.ProductCount,
		OrderComment:   res.OrderComment,
		OrderAddress:   res.OrderAddress,
		OrderStatus:    res.OrderStatus,
		IsCompleted:    res.IsCompleted,
	}
	return currentOrder, nil
}

func (ors *OrderService) AcceptOrder(ctx context.Context, orderID *OrderID) (*Response, error) {
	resp, err := ors.Database.AcceptOrder(&config.OrderID{Id: orderID.Id})
	if err != nil {
		return &Response{Code: 500, Message: "AcceptOrder service error"}, err
	}
	return &Response{Code: resp.Code, Message: resp.Message}, err
}

func (ors *OrderService) CompleteOrder(ctx context.Context, orderID *OrderID) (*Response, error) {
	resp, err := ors.Database.CompleteOrder(&config.OrderID{Id: orderID.Id})
	if err != nil {
		return &Response{Code: 500, Message: "CompleteOrder service error"}, err
	}
	return &Response{Code: resp.Code, Message: resp.Message}, err
}

func (ors *OrderService) CheckUserBlock(ctx context.Context, user *User) (*UserBlock, error) {
	if !validateUser(user) {
		return &UserBlock{}, wrongUserStruct
	}
	blockInfo, err := ors.Database.CheckUserBlock(&config.User{Username: user.Username})
	if err != nil {
		return &UserBlock{}, err
	}

	userBlock := &UserBlock{
		Id:                blockInfo.Id,
		IntruderUsername:  blockInfo.IntruderUsername,
		ModeratorUsername: blockInfo.ModeratorUsername,
		BanReason:         blockInfo.BanReason,
		BanDate:           blockInfo.BanDate,
		ExpiresAt:         blockInfo.ExpiresAt,
	}

	return userBlock, nil
}

func (ors *OrderService) BlockUser(ctx context.Context, user *UserBlock) (*Response, error) {
	if !validateUserBlock(user) {
		return &Response{Code: 400, Message: "Wrong request"}, wrongUserBlockStruct
	}

	userBlock := &config.UserBlock{
		Id:                user.Id,
		IntruderUsername:  user.IntruderUsername,
		ModeratorUsername: user.ModeratorUsername,
		BanReason:         user.BanReason,
		BanDate:           user.BanDate,
		ExpiresAt:         user.ExpiresAt,
	}

	resp, err := ors.Database.BlockUser(userBlock)
	if err != nil {
		return &Response{Code: 400, Message: "Internal server error"}, err
	}
	return &Response{Code: resp.Code, Message: resp.Message}, nil
}
func (ors *OrderService) UnblockUser(ctx context.Context, user *User) (*Response, error) {
	if !validateUser(user) {
		return &Response{Code: 400, Message: "Wrong request"}, wrongUserStruct
	}
	resp, err := ors.Database.UnblockUser(&config.User{Username: user.Username})
	if err != nil {
		return &Response{Code: 400, Message: "Internal server error"}, err
	}
	return &Response{Code: resp.Code, Message: resp.Message}, nil
}

// Validations for incoming structs

func validateOrder(order *Order) bool {
	if order.SellerUsername == "" || order.BuyerUsername == "" {
		return false
	} else if order.ProductId == 0 || order.ProductCount < 0 || order.ProductCount > 100 {
		return false
	} else if order.OrderAddress == "" {
		return false
	}
	return true
}

func validateUserBlock(userBlock *UserBlock) bool {
	if userBlock.IntruderUsername == "" || userBlock.ModeratorUsername == "" {
		return false
	} else if userBlock.BanReason == "" || userBlock.BanDate == "" || userBlock.ExpiresAt == "" {
		return false
	}
	return true
}

func validateUser(user *User) bool {
	if user.Username == "" {
		return false
	}
	return true
}

func validateBuyer(buyer *Buyer) bool {
	if buyer.BuyerUsername == "" {
		return false
	}
	return true
}
