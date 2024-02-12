package requestModels

type AdminSignInReq struct {
	Email    string `json:"Login_email" validate:"required,email"`
	Password string `json:"Login_password" validate:"required,gte=3"`
}
type BlockUserReq struct {
	Email string `json:"email" validate:"required,email"`
}

//unblock user request
type UnblockUserReq struct {
	Email string `json:"email" validate:"required,email"`
}

//block seller request
type BlockSellerReq struct {
	Email string `json:"email" validate:"required,email"`
}

//unblock seller request
type UnblockSellerReq struct {
	Email string `json:"email" validate:"required,email"`
}

//VerifySellerReq
type VerifySellerReq struct {
	SellerID  uint   `json:"seller_id" validate:"required,numeric"`
}
