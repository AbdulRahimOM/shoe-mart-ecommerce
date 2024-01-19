package requestModels

type MakeOrderReq struct {
	UserID            uint                `json:"userID" validate:"required,number"`
	AddressID         uint                `json:"addressID" validate:"required,number"`
	CouponID          uint                `json:"couponID" validate:"number"`
	PaymentMethod     string              `json:"paymentMethod" validate:"required,gte=2"`
	// ProductIDQuantity []MakeOrderProducts `json:"orderItems"`
}

type MakeOrderProducts struct {
	ProductID uint `json:"productID" validate:"required,number"`
	Quantity  uint `json:"quantity" validate:"required,number"`
}

//CancelOrderReq
type CancelOrderReq struct {
	OrderID uint `json:"orderID" validate:"required,number"`
}

//return order req
type ReturnOrderReq struct {
	OrderID uint `json:"orderID" validate:"required,number"`
}

//MarkOrderAsDeliveredReq
type MarkOrderAsDeliveredReq struct {
	OrderID uint `json:"orderID" validate:"required,number"`
}