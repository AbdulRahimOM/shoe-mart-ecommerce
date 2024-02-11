package response

import "MyShoo/internal/domain/entities"

type AdminLoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Error   string `json:"error"`
}
type GetUsersListResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error"`
	UsersList []entities.UserDetails `json:"usersList"`
}
type BlockUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
type UnblockUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
type GetSellersListResponse struct {
	Status      string                    `json:"status"`
	Message     string                    `json:"message"`
	Error       string                    `json:"error"`
	SellersList []entities.PwMaskedSeller `json:"sellersList"`
}

// BlockSellerResponse
type BlockSellerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// UnblockSellerResponse
type UnblockSellerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// GetCouponRes
type GetCouponRes struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Coupons []entities.Coupon `json:"coupons"`
}
