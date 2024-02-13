package repo

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type ICartRepo interface {
	//returns true if product exists in cart, quantity of product in cart, error
	DoProductExistAlready(cart *entities.Cart) (bool, uint, error)
	AddToCart(cart *entities.Cart) error
	GetCart(userID uint) (*[]entities.Cart, float32, error)
	DeleteFromCart(req *request.DeleteFromCartReq) error
	UpdateCartItemQuantity(cart *entities.Cart) error
	IsCartEmpty(userID uint) (bool, error)
	ClearCartOfUser(userID uint) error
	GetQuantityAndPriceOfCart(userID uint) (uint, float32, string, error)
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
	//order related_____________________________________________________________
	MakeOrder_UpdateStock_ClearCart(order *entities.Order, orderItems *[]entities.OrderItem) (uint, error)
	MakeOrder(order *entities.Order, orderItems *[]entities.OrderItem) (uint, error)
	UpdateOrderToPaid_UpdateStock_ClearCart(orderID uint) (*entities.Order, error)
	GetOrdersOfUser(userID uint, resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, error)
	GetOrders(resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, error)
	GetAllOrders() (*[]entities.Order, error)
	// GetOrderItemsByOrderID(orderID uint) (*[]entities.OrderItem, error)
	GetOrderItemsPQRByOrderID(orderID uint) (*[]response.PQMS, error)
	DoOrderExistByID(orderID uint) (bool, error)
	GetUserIDByOrderID(orderID uint) (uint, error)
	GetOrderStatusByID(orderID uint) (string, error)
	GetOrderSummaryByID(orderID uint) (*entities.Order, error)
	CancelOrder(orderID uint) error
	ReturnOrderRequest(orderID uint) error
	MarkOrderAsReturned(orderID uint) error

	//mark order as delivered and change payment-status to "paid" in case of COD
	MarkOrderAsDelivered(orderID uint) error

	//Online transaction related_______________________________________________________
	GetOrderByTransactionID(transactionID string) (uint, error)
	UpdateOrderTransactionID(orderID uint, transactionID string) error
	GetPaymentStatusByID(orderID uint) (string, error)

	//coupon related_______________________________________________________
	DoCouponExistByCode(code string) (bool, error)
	CreateNewCoupon(coupon *entities.Coupon) error
	BlockCoupon(couponID uint) error
	UnblockCoupon(couponID uint) error
	GetAllCoupons() (*[]entities.Coupon, string, error)
	GetExpiredCoupons() (*[]entities.Coupon, string, error)
	GetActiveCoupons() (*[]entities.Coupon, string, error)
	GetUpcomingCoupons() (*[]entities.Coupon, string, error)
	GetCouponByID(couponID uint) (*entities.Coupon, string, error)
	GetCouponUsageCount(userID uint, couponID uint) (uint, string, error)

	//upload related_______________________________________________________
	UploadInvoice(file string, fileName string) (string, error)
}
