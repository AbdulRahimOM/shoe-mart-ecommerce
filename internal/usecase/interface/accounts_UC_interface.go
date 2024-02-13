package usecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
)

type IAdminUC interface {
	//system related
	RestartConfig() error

	//self
	SignIn(req *request.AdminSignInReq) (*string, error)

	//user related
	GetUsersList() (*[]entities.UserDetails, error)
	BlockUser(req *request.BlockUserReq) error
	UnblockUser(req *request.UnblockUserReq) error

	//seller related
	GetSellersList() (*[]entities.PwMaskedSeller, error)
	BlockSeller(req *request.BlockSellerReq) error
	UnblockSeller(req *request.UnblockSellerReq) error
	VerifySeller(req *request.VerifySellerReq) error
}
type ISellerUC interface {
	SignUp(req *request.SellerSignUpReq) (*string, error)
	SignIn(req *request.SellerSignInReq) (*string, error)
}
type IUserUC interface {
	SignUp(req *request.UserSignUpReq) (*string, error)
	SignIn(req *request.UserSignInReq) (*string, error)
	SendOtp(phone string) error
	VerifyOtp(phone string, email string, otp string) (bool, error)

	//forgot password related
	GetUserByEmail(email string) (*entities.User, error)
	SendOtpForPWChange(*entities.User) (*string, error)
	VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, error)
	ResetPassword(id uint, newPassword *string) error

	//address related
	AddUserAddress(req *request.AddUserAddress) error
	EditUserAddress(req *request.EditUserAddress) error
	DeleteUserAddress(req *request.DeleteUserAddress) error
	GetUserAddresses(userID uint) (*[]entities.UserAddress, error)

	GetProfile(userID uint) (*entities.UserDetails, error)
	EditProfile(userID uint, req *request.EditProfileReq) error
}
