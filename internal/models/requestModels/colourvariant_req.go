package request

// colour variant related requests
type AddColourVariantReq struct {
	Colour    string  `json:"colour" validate:"required,gte=3" form:"colour"`
	ModelID   uint    `json:"modelId" validate:"required,number" form:"modelId"`
	MRP       float32 `json:"mrp" validate:"required,number" form:"mrp"`
	SalePrice float32 `json:"salePrice" validate:"required,number" form:"salePrice"`
	// ImageURL  string  `json:"-" form:"imageUrl"` //update required for: swagger ignore?
}

// edit colour variant
type EditColourVariantReq struct {
	ID        uint    `json:"id" validate:"required,number"`
	Colour    string  `json:"colour" validate:"required,gte=3"`
	ModelID   uint    `json:"modelId" validate:"required,number"`
	MRP       float32 `json:"mrp" validate:"required,number"`
	SalePrice float32 `json:"salePrice" validate:"required,number"`
}
