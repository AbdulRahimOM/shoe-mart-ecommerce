package usecase

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
)

type IAdminUC interface {
	//system related
	RestartConfig() *e.Error

	//self
	SignIn(req *request.AdminSignInReq) (*string, *e.Error)

	//user related
	GetUsersList() (*[]entities.UserDetails, *e.Error)
	BlockUser(req *request.BlockUserReq) *e.Error
	UnblockUser(req *request.UnblockUserReq) *e.Error

	//seller related
	GetSellersList() (*[]entities.PwMaskedSeller, *e.Error)
	BlockSeller(req *request.BlockSellerReq) *e.Error
	UnblockSeller(req *request.UnblockSellerReq) *e.Error
	VerifySeller(req *request.VerifySellerReq) *e.Error
}
type ISellerUC interface {
	SignUp(req *request.SellerSignUpReq) (*string, *e.Error)
	SignIn(req *request.SellerSignInReq) (*string, *e.Error)
}

// type IUserUC interface {
// 	SignUp(req *request.UserSignUpReq) (*string, *e.Error)
// 	SignIn(req *request.UserSignInReq) (*string, *e.Error)
// 	SendOtp(phone string) *e.Error
// 	VerifyOtp(phone string, email string, otp string) (bool, *e.Error)

// 	//forgot password related
// 	GetUserByEmail(email string) (*entities.User, *e.Error)
// 	SendOtpForPWChange(*entities.User) (*string, *e.Error)
// 	VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, *e.Error)
// 	ResetPassword(id uint, newPassword *string) *e.Error

// 	//address related
// 	AddUserAddress(req *request.AddUserAddress) *e.Error
// 	EditUserAddress(req *request.EditUserAddress) *e.Error
// 	DeleteUserAddress(req *request.DeleteUserAddress) *e.Error
// 	GetUserAddresses(userID uint) (*[]entities.UserAddress, *e.Error)

// 	GetProfile(userID uint) (*entities.UserDetails, *e.Error)
// 	EditProfile(userID uint, req *request.EditProfileReq) *e.Error
// }
