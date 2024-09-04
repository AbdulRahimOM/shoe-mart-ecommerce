package response

import "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"

type AdminLoginResponse struct {
	//Status  string `json:"status"`
	//Message string `json:"message"`
	Token string `json:"token"`
}

type GetUsersListResponse struct {
	//Status    string                 `json:"status"`
	//Message   string                 `json:"message"`
	UsersList []entities.UserDetails `json:"usersList"`
}

type GetSellersListResponse struct {
	//Status      string                    `json:"status"`
	//Message     string                    `json:"message"`
	SellersList []entities.PwMaskedSeller `json:"sellersList"`
}

// GetCouponRes
type GetCouponRes struct {
	Coupons []entities.Coupon `json:"coupons"`
}
