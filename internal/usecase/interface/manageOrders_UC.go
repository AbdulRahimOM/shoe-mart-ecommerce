package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type ICartUC interface {
	AddToCart(req *requestModels.AddToCartReq) error
	GetCart(userID uint) (*[]response.ResponseCartItems, float32, error)
	DeleteFromCart(req *requestModels.DeleteFromCartReq) error
	ClearCartOfUser(userID uint) error
}

type IWishListsUC interface {
	CreateWishList(userID uint, req *requestModels.CreateWishListReq) error
	//add to wishlist
	AddToWishList(userID uint, req *requestModels.AddToWishListReq) error
	//remove from wishlist
	RemoveFromWishList(userID uint, req *requestModels.RemoveFromWishListReq) error
	GetAllWishLists(userID uint) (*[]entities.WishList, int, error)
	GetWishListByID(userID uint, wishListID uint) (*string, *[]response.ResponseProduct2, int, error)
	//get wishlist of user
	// GetWishList(userID uint) (*[]response.ResponseWishList, error)
	// //get list  of all wishlists
	// GetAllWishLists() (*[]response.ResponseWishList, error)
}

type IOrderUC interface {
	//returns orderInfo, message, error
	MakeOrder(req *requestModels.MakeOrderReq) (*entities.OrderInfo, string, error)
	GetOrdersOfUser(userID uint, page int, limit int) (*[]response.ResponseOrderInfo, string, error)
	GetOrders(page int, limit int) (*[]response.ResponseOrderInfo, string, error)
	CancelOrderByUser(orderID uint, userID uint) (string, error)
	CancelOrderByAdmin(orderID uint) (string, error)
	ReturnOrderRequestByUser(orderID uint, userID uint) (string, error)
	MarkOrderAsReturned(orderID uint) (string, error)
	MarkOrderAsDelivered(orderID uint) (string, error)
}
