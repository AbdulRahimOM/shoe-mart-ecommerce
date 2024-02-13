package request

// wishlist
type CreateWishListReq struct {
	Name string `json:"name" validate:"required"`
}

type AddToWishListReq struct {
	WishListID uint `json:"wishListID" validate:"required,number"`
	ProductID  uint `json:"productID" validate:"required,number"`
}
type RemoveFromWishListReq struct {
	WishListID uint `json:"wishListID" validate:"required,number"`
	ProductID  uint `json:"productID" validate:"required,number"`
}
