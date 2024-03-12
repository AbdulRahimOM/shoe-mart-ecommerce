package usecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type ICartUC interface {
	AddToCart(req *request.AddToCartReq) *e.Error
	GetCart(userID uint) (*[]response.ResponseCartItems, float32, *e.Error)
	DeleteFromCart(req *request.DeleteFromCartReq) *e.Error
	ClearCartOfUser(userID uint) *e.Error
}

type IWishListsUC interface {
	CreateWishList(userID uint, req *request.CreateWishListReq) *e.Error
	//add to wishlist
	AddToWishList(userID uint, req *request.AddToWishListReq) *e.Error
	//remove from wishlist
	RemoveFromWishList(userID uint, req *request.RemoveFromWishListReq) *e.Error
	GetAllWishLists(userID uint) (*[]entities.WishList, int, *e.Error)
	GetWishListByID(userID uint, wishListID uint) (*string, *[]response.ResponseProduct2, int, *e.Error)
}

type IOrderUC interface {
	GetAddressForCheckout(userID uint) (*[]entities.UserAddress, uint, float32, *e.Error)
	SetAddressGetCoupons(userID uint, req *request.SetAddressForCheckOutReq) (*response.SetAddrGetCouponsResponse, *e.Error)
	SetCouponGetPaymentMethods(userID uint, req *request.SetCouponForCheckoutReq) (*response.GetPaymentMethodsForCheckoutResponse, *e.Error)

	//returns orderInfo, message, *e.Error
	MakeOrder(req *request.MakeOrderReq) (*entities.OrderInfo, *response.ProceedToPaymentInfo, *e.Error)
	GetOrdersOfUser(userID uint, page int, limit int) (*[]response.ResponseOrderInfo,  *e.Error)
	GetOrders(page int, limit int) (*[]response.ResponseOrderInfo, *e.Error)
	CancelOrderByUser(orderID uint, userID uint) *e.Error
	CancelOrderByAdmin(orderID uint) *e.Error
	ReturnOrderRequestByUser(orderID uint, userID uint) *e.Error
	MarkOrderAsReturned(orderID uint) *e.Error
	MarkOrderAsDelivered(orderID uint) *e.Error
	// ProceedToPayment(req *request.ProceedToPaymentReq) (*response.ProceedToPaymentInfo,string, *e.Error)

	GetInvoiceOfOrder(userID uint, orderID uint) (*string, *e.Error)

	CreateNewCoupon(req *request.NewCouponReq) *e.Error
	BlockCoupon(req *request.BlockCouponReq) *e.Error
	UnblockCoupon(req *request.UnblockCouponReq) *e.Error
	GetAllCoupons() (*[]entities.Coupon, *e.Error)
	GetExpiredCoupons() (*[]entities.Coupon, *e.Error)
	GetActiveCoupons() (*[]entities.Coupon, *e.Error)
	GetUpcomingCoupons() (*[]entities.Coupon, *e.Error)
}
