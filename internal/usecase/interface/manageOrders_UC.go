package usecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type ICartUC interface {
	AddToCart(req *request.AddToCartReq) error
	GetCart(userID uint) (*[]response.ResponseCartItems, float32, error)
	DeleteFromCart(req *request.DeleteFromCartReq) error
	ClearCartOfUser(userID uint) error
}

type IWishListsUC interface {
	CreateWishList(userID uint, req *request.CreateWishListReq) error
	//add to wishlist
	AddToWishList(userID uint, req *request.AddToWishListReq) error
	//remove from wishlist
	RemoveFromWishList(userID uint, req *request.RemoveFromWishListReq) error
	GetAllWishLists(userID uint) (*[]entities.WishList, int, error)
	GetWishListByID(userID uint, wishListID uint) (*string, *[]response.ResponseProduct2, int, error)
}

type IOrderUC interface {
	GetAddressForCheckout(userID uint) (*[]entities.UserAddress, uint, float32, string, error)
	SetAddressGetCoupons(userID uint, req *request.SetAddressForCheckOutReq) (*response.SetAddrGetCouponsResponse, string, error)
	SetCouponGetPaymentMethods(userID uint, req *request.SetCouponForCheckoutReq) (*response.GetPaymentMethodsForCheckoutResponse, string, error)

	//returns orderInfo, message, error
	MakeOrder(req *request.MakeOrderReq) (*entities.OrderInfo, *response.ProceedToPaymentInfo, string, error)
	GetOrdersOfUser(userID uint, page int, limit int) (*[]response.ResponseOrderInfo, string, error)
	GetOrders(page int, limit int) (*[]response.ResponseOrderInfo, string, error)
	CancelOrderByUser(orderID uint, userID uint) (string, error)
	CancelOrderByAdmin(orderID uint) (string, error)
	ReturnOrderRequestByUser(orderID uint, userID uint) (string, error)
	MarkOrderAsReturned(orderID uint) (string, error)
	MarkOrderAsDelivered(orderID uint) (string, error)
	// ProceedToPayment(req *request.ProceedToPaymentReq) (*response.ProceedToPaymentInfo,string, error)

	GetInvoiceOfOrder(userID uint, orderID uint) (*string, string, error)

	CreateNewCoupon(req *request.NewCouponReq) (string, error)
	BlockCoupon(req *request.BlockCouponReq) (string, error)
	UnblockCoupon(req *request.UnblockCouponReq) (string, error)
	GetAllCoupons() (*[]entities.Coupon, string, error)
	GetExpiredCoupons() (*[]entities.Coupon, string, error)
	GetActiveCoupons() (*[]entities.Coupon, string, error)
	GetUpcomingCoupons() (*[]entities.Coupon, string, error)
}
