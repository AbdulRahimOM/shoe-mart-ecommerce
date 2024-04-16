package request

type AdminSignInReq struct {
	Email    string `json:"Login_email" validate:"required,email"`
	Password string `json:"Login_password" validate:"required,gte=3"`
}
type BlockUserReq struct {
	UserID uint `json:"user_id" validate:"required,numeric"`
}

// unblock user request
type UnblockUserReq struct {
	UserID uint `json:"user_id" validate:"required,numeric"`
}

// block seller request
type BlockSellerReq struct {
	SellerID uint `json:"seller_id" validate:"required,numeric"`
}

// unblock seller request
type UnblockSellerReq struct {
	SellerID uint `json:"seller_id" validate:"required,numeric"`
}

// VerifySellerReq
type VerifySellerReq struct {
	SellerID uint `json:"seller_id" validate:"required,numeric"`
}
