package usecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
)

type IUserUC interface {
	SignUp(req *request.UserSignUpReq) (*string, *e.Error)
	SignIn(req *request.UserSignInReq) (*string, *e.Error)
	SendOtp(phone string) *e.Error
	VerifyOtp(phone string, email string, otp string) (bool, *e.Error)

	//forgot password related
	GetUserByEmail(email string) (*entities.User, *e.Error)
	SendOtpForPWChange(*entities.User) (*string, *e.Error)
	VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, *e.Error)
	ResetPassword(id uint, newPassword *string) *e.Error

	//address related
	AddUserAddress(req *request.AddUserAddress) *e.Error
	EditUserAddress(req *request.EditUserAddress) *e.Error
	DeleteUserAddress(req *request.DeleteUserAddress) *e.Error
	GetUserAddresses(userID uint) (*[]entities.UserAddress, *e.Error)

	GetProfile(userID uint) (*entities.UserDetails, *e.Error)
	EditProfile(userID uint, req *request.EditProfileReq) *e.Error
}
