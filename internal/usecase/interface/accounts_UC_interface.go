package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
)

type IAdminUC interface {
	//system related
	RestartConfig() error

	//self
	SignIn(req *requestModels.AdminSignInReq) (*string, error)

	//user related
	GetUsersList() (*[]entities.UserDetails, error)
	BlockUser(req *requestModels.BlockUserReq) error
	UnblockUser(req *requestModels.UnblockUserReq) error

	//seller related
	GetSellersList() (*[]entities.PwMaskedSeller, error)
	BlockSeller(req *requestModels.BlockSellerReq) error
	UnblockSeller(req *requestModels.UnblockSellerReq) error
}
type ISellerUC interface {
	SignUp(req *requestModels.SellerSignUpReq) (*string, error)
	SignIn(req *requestModels.SellerSignInReq) (*string, error)
}
type IUserUC interface {
	SignUp(req *requestModels.UserSignUpReq) (*string, error)
	SignIn(req *requestModels.UserSignInReq) (*string, error)
	SendOtp(phone string) error
	VerifyOtp(phone string, email string, otp string) (bool, error)

	//forgot password related
	GetUserByEmail(email string) (*entities.User, error)
	SendOtpForPWChange(*entities.User) (*string, error)
	VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, error)
	ResetPassword(id uint, newPassword *string) error

	//address related
	AddUserAddress(req *requestModels.AddUserAddress) error
	EditUserAddress(req *requestModels.EditUserAddress) error
	DeleteUserAddress(req *requestModels.DeleteUserAddress) error
	GetUserAddresses(userID uint) (*[]entities.UserAddress, error)

	GetProfile(userID uint) (*entities.UserDetails, error)
	EditProfile(userID uint, req *requestModels.EditProfileReq) error
}
