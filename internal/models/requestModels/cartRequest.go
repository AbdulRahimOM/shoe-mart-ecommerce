package requestModels

// cart
type AddToCartReq struct {
	UserID    uint `json:"userID" validate:"required,number"`
	ProductID uint `json:"productID" validate:"required,number"`
}
type DeleteFromCartReq struct {
	UserID    uint `json:"userID" validate:"required,number"`
	ProductID uint `json:"productID" validate:"required,number"`
}
type ProductIDQuantity struct {
	ProductID uint `json:"productID" validate:"required,number"`
	Quantity  uint `json:"quantity" validate:"required,number"`
}
type IncreaseQuantityReq struct {
	UserID            uint                `json:"userID" validate:"required,number"`
	ProductIDQuantity []ProductIDQuantity `json:"productIDQuantity" validate:"required"`
}
type EditCartReq struct {
	UserID    uint `json:"userID" validate:"required,number"`
	ProductID uint `json:"productID" validate:"required,number"`
	Quantity  uint `json:"quantity" validate:"required,number"`
}
