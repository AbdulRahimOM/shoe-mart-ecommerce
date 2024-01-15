package requestModels

// Category related requests
type AddCategoryReq struct {
	Name string `json:"name" validate:"required,gte=3"`
}
type EditCategoryReq struct {
	OldName string `json:"oldName" validate:"required,gte=3"`
	// ID uint `json:"id" validate:"required,number"`
	NewName string `json:"newName" validate:"required,gte=3"`
}
type DeleteCategoryReq struct {
	Name string `json:"name" validate:"required,gte=3"`
	// ID uint `json:"id" validate:"required,number"`
}