package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
)

type ICartRepo interface {
	//returns true if product exists in cart, quantity of product in cart, error
	DoProductExistAlready(cart *entities.Cart) (bool, uint, error)
	AddToCart(cart *entities.Cart) error
	GetCart(userID uint) (*[]entities.Cart, error)
	DeleteFromCart(req *requestModels.DeleteFromCartReq) error
	UpdateCartItemQuantity(cart *entities.Cart) error
	IsCartEmpty(userID uint) (bool, error)
	ClearCartOfUser(userID uint) error
}

type IWishListsRepo interface {
	DoesWishListExistWithName(userID uint, name string) (bool, error)
	CreateWishList(userID uint, wishList *entities.WishList) error
	DoesThisWishListExistForUser(userID uint, wishListID uint) (bool, error)
	IsProductInWishList(productID uint, wishListID uint) (bool, error)
	AddToWishList(productID uint, wishListID uint) error
	RemoveFromWishList(productID uint, wishListID uint) error
	GetAllWishLists(userID uint) (*[]entities.WishList, error)
	GetWishListByID(userID uint, wishListID uint) (*string, *[]entities.Product, error)
}

type IOrderRepo interface {
	MakeOrder_UpdateStock_ClearCart(order *entities.Order, orderItems *[]entities.OrderItem) (*entities.Order, error)

	MakeOrder(order *entities.Order, orderItems *[]entities.OrderItem) (*entities.Order, error)
	UpdateOrderToPaid_UpdateStock_ClearCart(orderID uint) (*entities.Order, error)

	GetOrdersOfUser(userID uint, resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, error)
	GetOrders(resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, error)
	GetAllOrders() (*[]entities.Order, error)
	DoOrderExistByID(orderID uint) (bool, error)
	GetUserIDByOrderID(orderID uint) (uint, error)
	GetOrderStatusByID(orderID uint) (string, error)
	GetOrderSummaryByID(orderID uint) (*entities.Order, error)
	CancelOrder(orderID uint) error
	ReturnOrderRequest(orderID uint) error
	MarkOrderAsReturned(orderID uint) error

	//mark order as delivered and change payment-status to "paid" in case of COD
	MarkOrderAsDelivered(orderID uint) error

	GetOrderByTransactionID(transactionID string) (uint, error)
	UpdateOrderTransactionID(orderID uint, transactionID string) error
}
