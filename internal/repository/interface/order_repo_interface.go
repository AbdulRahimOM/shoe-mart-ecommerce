package repo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type ICartRepo interface {
	//returns true if product exists in cart, quantity of product in cart, *e.Error
	DoProductExistAlready(cart *entities.Cart) (bool, uint, *e.Error)
	AddToCart(cart *entities.Cart) *e.Error
	GetCart(userID uint) (*[]entities.Cart, float32, *e.Error)
	DeleteFromCart(req *request.DeleteFromCartReq) *e.Error
	UpdateCartItemQuantity(cart *entities.Cart) *e.Error
	IsCartEmpty(userID uint) (bool, *e.Error)
	ClearCartOfUser(userID uint) *e.Error
	GetQuantityAndPriceOfCart(userID uint) (uint, float32, *e.Error)
}

type IWishListsRepo interface {
	DoesWishListExistWithName(userID uint, name string) (bool, *e.Error)
	CreateWishList(userID uint, wishList *entities.WishList) *e.Error
	GetUserIDOfWishList(wishListID uint) (uint, *e.Error)
	IsProductInWishList(productID uint, wishListID uint) (bool, *e.Error)
	AddToWishList(productID uint, wishListID uint) *e.Error
	RemoveFromWishList(productID uint, wishListID uint) *e.Error
	GetAllWishLists(userID uint) (*[]entities.WishList, *e.Error)
	GetWishListByID(userID uint, wishListID uint) (*string, *[]entities.Product, *e.Error)
}

type IOrderRepo interface {
	//order related_____________________________________________________________
	MakeOrder_UpdateStock_ClearCart(order *entities.Order, orderItems *[]entities.OrderItem) (uint, *e.Error)
	MakeOrder(order *entities.Order, orderItems *[]entities.OrderItem) (uint, *e.Error)
	UpdateOrderToPaid_UpdateStock_ClearCart(orderID uint) (*entities.Order, *e.Error)
	GetOrdersOfUser(userID uint, resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, *e.Error)
	GetOrders(resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, *e.Error)
	GetAllOrders() (*[]entities.Order, *e.Error)
	GetOrderItemsPQRByOrderID(orderID uint) (*[]response.PQMS, *e.Error)
	GetUserIDByOrderID(orderID uint) (uint, *e.Error)
	GetOrderStatusByID(orderID uint) (string, *e.Error)
	GetOrderSummaryByID(orderID uint) (*entities.Order, *e.Error)
	CancelOrder(orderID uint) *e.Error
	ReturnOrderRequest(orderID uint) *e.Error
	MarkOrderAsReturned(orderID uint) *e.Error

	//mark order as delivered and change payment-status to "paid" in case of COD
	MarkOrderAsDelivered(orderID uint) *e.Error

	//Online transaction related_______________________________________________________
	GetOrderByTransactionID(transactionID string) (uint, *e.Error)
	UpdateOrderTransactionID(orderID uint, transactionID string) *e.Error
	GetPaymentStatusByID(orderID uint) (string, *e.Error)

	//coupon related_______________________________________________________
	IsCouponCodeTaken(code string) (bool, *e.Error)
	CreateNewCoupon(coupon *entities.Coupon) *e.Error
	BlockCoupon(couponID uint) *e.Error
	UnblockCoupon(couponID uint) *e.Error
	GetAllCoupons() (*[]entities.Coupon, *e.Error)
	GetExpiredCoupons() (*[]entities.Coupon, *e.Error)
	GetActiveCoupons() (*[]entities.Coupon, *e.Error)
	GetUpcomingCoupons() (*[]entities.Coupon, *e.Error)
	GetCouponByID(couponID uint) (*entities.Coupon, *e.Error)
	GetCouponUsageCount(userID uint, couponID uint) (uint, *e.Error)

	//upload related_______________________________________________________
	UploadInvoice(file string, fileName string) (string, *e.Error)
}
