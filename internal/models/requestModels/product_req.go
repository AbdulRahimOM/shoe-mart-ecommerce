package request

type AddStockReq struct {
	SellerID         uint `json:"sellerId" validate:"required,number"`
	ProductID        uint `json:"productId" validate:"required,number"`
	AddingStockCount uint `json:"addingStockCount" validate:"required,number"`
}
type EditStockReq struct {
	SellerID          uint `json:"sellerId" validate:"required,number"`
	ProductID         uint `json:"productId" validate:"required,number"`
	UpdatedStockCount uint `json:"updatedStockCount" validate:"required,number"`
}
