package response

import "MyShoo/internal/domain/entities"

type AdminLoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}
type GetUsersListResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	UsersList []entities.UserDetails `json:"usersList"`
}
// type BlockUserResponse struct {
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// }
// type UnblockUserResponse struct {
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// }
type GetSellersListResponse struct {
	Status      string                    `json:"status"`
	Message     string                    `json:"message"`
	SellersList []entities.PwMaskedSeller `json:"sellersList"`
}

// // BlockSellerResponse
// type BlockSellerResponse struct {
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// }

// // UnblockSellerResponse
// type UnblockSellerResponse struct {
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// }

// GetCouponRes
type GetCouponRes struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Coupons []entities.Coupon `json:"coupons"`
}
