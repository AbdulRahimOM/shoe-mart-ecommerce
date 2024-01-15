package requestModels

// Product related requests
type AddProductReq struct {
	SKUCode                   string `json:"skuCode"` //validate:"required,gte=3"
	SizeIndex                 uint   `json:"sizeIndex" validate:"required,number"`
	ColourVariantID           uint   `json:"colourVariantId" validate:"required,number"`
	DimensionalVariationIndex uint   `json:"dimensionalVariationIndex" validate:"required,number"`
	Stock                     uint   `json:"stock" validate:"required,number"`
}
type EditProductReq struct {
	ID                        uint   `json:"id" validate:"required,number"`
	SKUCode                   string `json:"skuCode"` // validate:"required,gte=3"`
	SizeIndex                 uint   `json:"sizeIndex" validate:"required,number"`
	ColourVariantID           uint   `json:"colourVariantId" validate:"required,number"`
	DimensionalVariationIndex uint   `json:"dimensionalVariationIndex" validate:"required,number"`
	Stock                     uint   `json:"stock" validate:"required,number"`
}
type DeleteProductReq struct {
	ID uint `json:"id" validate:"required,number"`
}

type AddStockReq struct{
	SellerID uint `json:"sellerId" validate:"required,number"`
	ProductID uint `json:"productId" validate:"required,number"`
	AddingStockCount uint `json:"addingStockCount" validate:"required,number"`
}
type EditStockReq struct{
	SellerID uint `json:"sellerId" validate:"required,number"`
	ProductID uint `json:"productId" validate:"required,number"`
	UpdatedStockCount uint `json:"updatedStockCount" validate:"required,number"`
}
