package request

// Model related requests
type AddModelReq struct {
	Name       string `json:"name" validate:"required,gte=3"`
	BrandID    uint   `json:"brandId" validate:"required,number"`
	CategoryID uint   `json:"categoryId" validate:"required,number"`
}

type EditModelReq struct {
	ID         uint   `json:"id" validate:"required,number"`
	Name       string `json:"name" validate:"required,gte=3"`
	BrandID    uint   `json:"brandId" validate:"required,number"`
	CategoryID uint   `json:"categoryId" validate:"required,number"`
}

// GetModelsUnderCategoryReq
type GetModelsUnderCategoryReq struct {
	CategoryID uint `json:"CategoryID" validate:"required,number"`
}
