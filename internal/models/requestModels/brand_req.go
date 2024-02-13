package request

// Brand related requests
type AddBrandReq struct {
	Name     string `json:"name" validate:"required,gte=3"`
	SellerID uint   `json:"sellerId" validate:"required,number"`
}
type EditBrandReq struct {
	// ID uint `json:"id" validate:"required,number"`
	OldName  string `json:"oldName" validate:"required,gte=3"`
	NewName  string `json:"newName" validate:"required,gte=3"`
	SellerID uint   `json:"sellerId" validate:"required,number"`
}
type DeleteBrandReq struct {
	// ID uint `json:"id" validate:"required,number"`
	Name string `json:"name" validate:"required,gte=3"`
}
