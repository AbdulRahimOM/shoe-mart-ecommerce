package requestModels

// colour variant related requests
type AddColourVariantReq struct {
	Colour    string `json:"colour" validate:"required,gte=3"`
	ModelID   uint   `json:"modelId" validate:"required,number"`
	MRP       float32   `json:"mrp" validate:"required,number"`
	SalePrice float32   `json:"salePrice" validate:"required,number"`
}

// edit colour variant
type EditColourVariantReq struct {
	ID        uint   `json:"id" validate:"required,number"`
	Colour    string `json:"colour" validate:"required,gte=3"`
	ModelID   uint   `json:"modelId" validate:"required,number"`
	MRP       float32   `json:"mrp" validate:"required,number"`
	SalePrice float32   `json:"salePrice" validate:"required,number"`
}
type DeleteColourVariantReq struct {
	ID uint `json:"id" validate:"required,number"`
}

//GetColourVariantsUnderModelReq
type GetColourVariantsUnderModelReq struct {
	ModelID uint `json:"modelId" validate:"required,number"`
}