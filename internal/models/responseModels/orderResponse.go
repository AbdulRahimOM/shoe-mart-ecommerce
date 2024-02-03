package response

import "MyShoo/internal/domain/entities"

// OrderResponse
type CODOrderResponse struct {
	Status    string             `json:"status"`
	Message   string             `json:"message"`
	Error     string             `json:"error"`
	OrderInfo entities.OrderInfo `json:"orderInfo"`
}

type OnlinePaymentOrderResponse struct {
	Status               string               `json:"status"`
	Message              string               `json:"message"`
	Error                string               `json:"error"`
	OrderInfo            entities.OrderInfo   `json:"orderInfo"`
	ProceedToPaymentInfo ProceedToPaymentInfo `json:"proceedToPaymentInfo"`
}
type PaidOrderResponse struct {
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	OrderInfo entities.Order `json:"order"`
}

// responseOrderInfo
type ResponseOrderInfo struct {
	OrderDetails entities.Order `json:"orderDetails"`
	OrderItems   []PQR          `json:"orderItems"`
}

type GetOrdersResponse struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Error      string              `json:"error"`
	OrdersInfo []ResponseOrderInfo `json:"ordersInfo"`
}

// get cart response
type GetCartResponse struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Error      string              `json:"error"`
	TotalValue float32             `json:"totalValue"`
	Cart       []ResponseCartItems `json:"cart"`
}

// responseCart
type ResponseCartItems struct {
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`

	FkProduct struct {
		SKUCode string `json:"skuCode"`
		Name    string `json:"name"`
		Stock   uint   `json:"stock"`
	} `json:"fkProduct"`
}

type GetAllWishListsResponse struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Error      string              `json:"error"`
	WishLists  []entities.WishList `json:"wishLists"`
	TotalCount int                 `json:"totalCount"`
}

type GetWishListByIDResponse struct {
	Status       string             `json:"status"`
	Message      string             `json:"message"`
	Error        string             `json:"error"`
	WishListName string             `json:"wishListName"`
	WishItems    []ResponseProduct2 `json:"wishItems"`
	TotalCount   int                `json:"totalCount"`
}

// responseProduct
type ResponseProduct2 struct {
	ID      uint   `gorm:"column:id;autoIncrement;primaryKey"`
	SKUCode string `gorm:"column:skuCode"`
	Name    string `gorm:"column:name;notNull"`
	// SizeIndex              uint    `gorm:"column:sizeIndex;notNull"`
	// DimensionalVariationID uint    `gorm:"column:dimensionalVariationID;notNull"`
	// Stock                  uint    `gorm:"column:stock;notNull"`
	MRP       float32 `gorm:"column:mrp;notNull"`
	SalePrice float32 `gorm:"column:salePrice;notNull"`
}

// GetInvoiceResponse
type GetInvoiceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Invoice string `json:"invoice"`
}

// checkOutInfo
type CheckOutInfo struct {
	ItemCount  uint                   `json:"itemsCount"`
	TotalValue float32                `json:"totalValue"`
	Addresses  []entities.UserAddress `json:"addresses"`
	Coupons    []entities.Coupon      `json:"coupons"`
}

// GetCheckoutResponse
type GetCheckoutResponse struct {
	Status       string       `json:"status"`
	Message      string       `json:"message"`
	CheckOutInfo CheckOutInfo `json:"checkOutInfo"`
}

// CheckoutEstimateResponse
type CheckoutEstimateResponse struct {
	Status         string   `json:"status"`
	Message        string   `json:"message"`
	ProductsValue  float32  `json:"productsValue"`
	Discount       float32  `json:"discount"`
	ShippingCharge float32  `json:"shippingCharge"`
	GrandTotal     float32  `json:"grandTotal"`
	PaymentMethods []string `json:"paymentMethods"`
}

type EstimateInfo struct {
	TotalValue float32 `json:"totalValue"`
	Discount   float32 `json:"discount"`
	Shipping   float32 `json:"shipping"`
	GrandTotal float32 `json:"grandTotal"`
}

// GetAddressesForCheckoutResponse
type GetAddressesForCheckoutResponse struct {
	Status       string                 `json:"status"`
	Message      string                 `json:"message"`
	Addresses    []entities.UserAddress `json:"addresses"`
	TotalQuantiy uint                   `json:"totalQuantiy"`
	TotalValue   float32                `json:"totalValue"`
}

// SetAddrGetCouponsResponse
type SetAddrGetCouponsResponse struct {
	Status       string               `json:"status"`
	Message      string               `json:"message"`
	Coupons      []ResponseCoupon    `json:"coupons"`
	Address      entities.UserAddress `json:"address"`
	TotalQuantiy uint                 `json:"totalQuantiy"`
	BillSumary   BillBeforeCoupon     `json:"billSummary"`
}

type BillBeforeCoupon struct {
	TotalProductsValue float32 `json:"totalProductsValue"`
	ShippingCharge     float32 `json:"shippingCharge"`
	GrandTotal         float32 `json:"grandTotal"`
}

// GetPaymentMethodsForCheckoutResponse
type GetPaymentMethodsForCheckoutResponse struct {
	Status              string               `json:"status"`
	Message             string               `json:"message"`
	Address             entities.UserAddress `json:"address"`
	TotalQuantiy        uint                 `json:"totalQuantiy"`
	BillSumary          BillAfterCoupon      `json:"billSummary"`
	PaymentMethods      []string             `json:"paymentMethods"`
	CODAvailability     bool                 `json:"codAvailability"`
	CODAvailabilityNote string               `json:"codAvailabilityNote"`
	WalletBalance       float32              `json:"walletBalance"`
}

type BillAfterCoupon struct {
	TotalProductsValue float32         `json:"totalProductsValue"`
	CouponApplied      bool            `json:"couponApplied"`
	Coupon             ResponseCoupon `json:"coupon"`
	CouponDiscount     float32         `json:"couponDiscount"`
	ShippingCharge     float32         `json:"shippingCharge"`
	GrandTotal         float32         `json:"grandTotal"`
}

type ResponseCoupon struct {
	ID            uint    `json:"id"`
	Code          string  `json:"code"`
	Type          string  `json:"type"`
	MinOrderValue float32 `json:"minOrderValue"`
	MaxDiscount   float32 `json:"maxDiscount"`
	Percentage    float32 `json:"percentage"`
	Description   string  `json:"description"`
}
