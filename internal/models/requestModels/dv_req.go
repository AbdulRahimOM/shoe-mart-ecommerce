package request

// Dimensional Variant
type AddDimensionalVariantReq struct {
	ColourVariantID uint `json:"colourVariantId" validate:"required,number"`
	DVIndex         uint `json:"dvIndex" validate:"required,number"`
}
type EditDimensionalVariantReq struct {
	ID              uint `json:"id" validate:"required,number"`
	ColourVariantID uint `json:"colourVariantId" validate:"required,number"`
	DVIndex         uint `json:"dvIndex" validate:"required,number"`
}
