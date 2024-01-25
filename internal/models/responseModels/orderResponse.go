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
	Status    string             `json:"status"`
	Message   string             `json:"message"`
	Error     string             `json:"error"`
	OrderInfo entities.OrderInfo `json:"orderInfo"`
	ProceedToPaymentInfo ProceedToPaymentInfo `json:"proceedToPaymentInfo"`
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
	Status       string `json:"status"`
	Message      string `json:"message"`
	Error        string `json:"error"`
	WishListName string `json:"wishListName"`
	WishItems    []ResponseProduct2 `json:"wishItems"`
	TotalCount int `json:"totalCount"`
}

// responseProduct
type ResponseProduct2 struct {
	ID                     uint    `gorm:"column:id;autoIncrement;primaryKey"`
	SKUCode                string  `gorm:"column:skuCode"`
	Name                   string  `gorm:"column:name;notNull"`
	// SizeIndex              uint    `gorm:"column:sizeIndex;notNull"`
	// DimensionalVariationID uint    `gorm:"column:dimensionalVariationID;notNull"`
	// Stock                  uint    `gorm:"column:stock;notNull"`
	MRP                    float32 `gorm:"column:mrp;notNull"`
	SalePrice              float32 `gorm:"column:salePrice;notNull"`
}
