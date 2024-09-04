package accountsusecase

import (
	"errors"
	"time"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
	hashpassword "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/hashPassword"
	jwttoken "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/jwt"
	otpManager "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/twilio"

	"github.com/jinzhu/copier"
)

var (
	errAddressNameAlreadyExists_409   = e.Error{Status: "failed", Msg: "address name already exists", Err: nil, StatusCode: 409}
	errAddressDoesNotBelongToUser_401 = e.Error{Status: "failed", Msg: "Corrupt request", Err: errors.New("address does not belong to user"), StatusCode: 401}
	errEmailIsAlreadyTaken_409        = e.Error{Status: "failed", Msg: "email already exists for another user", Err: nil, StatusCode: 409}
)

type UserUseCase struct {
	userRepo repoInterface.IUserRepo
}

func NewUserUseCase(repo repoInterface.IUserRepo) usecase.IUserUC {
	return &UserUseCase{userRepo: repo}
}

func (uc *UserUseCase) SendOtp(phone string) *e.Error {
	err := otpManager.SendOtp(phone)
	if err != nil {
		return e.SetError("Error while sending otp", err, 500)
	}
	return nil
}

func (uc *UserUseCase) VerifyOtp(phone string, email string, otp string) (bool, string, *e.Error) {
	if matchStatus, err := otpManager.VerifyOtp(phone, otp); err != nil {
		return false, "", e.SetError("Error while verifying otp", err, 500)
	} else if !matchStatus {
		return false, "", nil
	}

	//otp matched!!!, update status to verified
	err := uc.userRepo.UpdateUserStatus(email, "verified")
	if err != nil {
		return true, "", err
	}

	//get new token
	user, err := uc.userRepo.GetUserDetailsByEmail(email)
	if err != nil {
		return true, "", err
	}

	tokenString, errr := jwttoken.GenerateToken("password-to-be-set-user", user, time.Hour*24*30)
	if errr != nil {
		return true, "", e.SetError("error generating token", errr, 500)
	}

	return true, tokenString, uc.userRepo.UpdateUserStatus(email, "verified")

}
func (uc *UserUseCase) SignIn(req *request.UserSignInReq) (*string, *e.Error) {
	isEmailRegistered, err := uc.userRepo.IsEmailRegistered(req.Email)
	if err != nil {
		return nil, err
	}
	if !(isEmailRegistered) {
		return nil, errEmailNotRegistered_401
	}

	//get userpassword from database
	userForToken, err := uc.userRepo.GetPasswordAndUserDetailsByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//check for password
	if err := hashpassword.CompareHashedPassword(userForToken.Password, req.Password); err != nil {
		return nil, e.SetError("wrong password", err, 401)
	}

	//generate token
	tokenString, errr := jwttoken.GenerateToken("user", userForToken, time.Hour*24*30)
	if errr != nil {
		return nil, e.SetError("error generating token", errr, 500)
	}

	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}

func (uc *UserUseCase) SignUp(req *request.UserSignUpReq) (*string, *e.Error) {
	emailAlreadyUsed, err := uc.userRepo.IsEmailRegistered(req.Email)
	if err != nil {
		return nil, err
	}
	if emailAlreadyUsed {
		return nil, errEmailAlreadyUsed_401
	}

	hashedPwd, errr := hashpassword.Hashpassword(req.Password)
	if errr != nil {
		return nil, e.SetError("error while hashing pw", errr, 500)
	}
	var signingUser = entities.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPwd,
		Status:    "not verified",
		CreatedAt: time.Now(),
	}

	err = uc.userRepo.CreateUser(&signingUser)
	if err != nil {
		return nil, err
	}

	var userForToken = entities.UserDetails{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "not verified",
	}

	//generate token
	tokenString, errr := jwttoken.GenerateToken("user", userForToken, time.Hour*5)
	if errr != nil {
		return nil, e.SetError("error generating token", errr, 500)
	}
	return &tokenString, nil
}

func (uc *UserUseCase) AddUserAddress(req *request.AddUserAddress) *e.Error {
	//check if address name already exists
	doAddressNameAlreadyExists, err := uc.userRepo.DoAddressNameExists(req.AddressName)
	if err != nil {
		return err
	}
	if doAddressNameAlreadyExists {
		return &errAddressNameAlreadyExists_409
	}
	var address entities.UserAddress
	if err := copier.Copy(&address, &req); err != nil {
		return e.SetError("Error while copying request to cart entity", err, 500)
	}
	//add address
	return uc.userRepo.AddUserAddress(&address)
}

// EditUserAddress
func (uc *UserUseCase) EditUserAddress(userID uint, req *request.EditUserAddress) *e.Error {
	//check if userid is conforming
	userIDInAddress, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		return err
	}
	if userID != userIDInAddress {
		return &errAddressDoesNotBelongToUser_401
	}

	//get old address name
	oldAddressName, err := uc.userRepo.GetAddressNameByID(req.ID)
	if err != nil {
		return err
	}
	if *oldAddressName != req.AddressName {
		//check if AddressNameExistsForAnotherAddressOfUser,
		doAddressNameAlreadyExists, err := uc.userRepo.DoAddressNameExists(req.AddressName)
		if err != nil {
			return err
		}
		if doAddressNameAlreadyExists {
			return &errAddressNameAlreadyExists_409
		}
	}

	var address entities.UserAddress
	if err := copier.Copy(&address, &req); err != nil {
		return e.SetError("Error while copying request to cart entity", err, 500)
	}
	//edit address
	return uc.userRepo.EditUserAddress(&address)
}

// DeleteUserAddress
func (uc *UserUseCase) DeleteUserAddress(userID uint, req *request.DeleteUserAddress) *e.Error {
	//check if userid is conforming
	userIDInAddress, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		return err
	}
	if userID != userIDInAddress {
		return &errAddressDoesNotBelongToUser_401
	}

	//delete address
	return uc.userRepo.DeleteUserAddress(req.ID)
}

// GetUserAddresses
func (uc *UserUseCase) GetUserAddresses(userID uint) (*[]entities.UserAddress, *e.Error) {
	return uc.userRepo.GetUserAddresses(userID)
}

// GetProfile
func (uc *UserUseCase) GetProfile(userID uint) (*entities.UserDetails, *e.Error) {
	return uc.userRepo.GetProfile(userID)
}

// EditProfile
func (uc *UserUseCase) EditProfile(userID uint, req *request.EditProfileReq) *e.Error {
	//check if email is not changed
	oldEmail, err := uc.userRepo.GetEmailByUserID(userID)
	if err != nil {
		return err
	}
	if *oldEmail != req.Email {
		//check if email already exists
		doEmailAlreadyExists, err := uc.userRepo.IsEmailRegistered(req.Email)
		if err != nil {
			return err
		}
		if doEmailAlreadyExists {
			return &errEmailIsAlreadyTaken_409
		}
	}

	//edit profile
	return uc.userRepo.EditProfile(userID, req)
}

// GetUserPhoneByEmail
func (uc *UserUseCase) GetUserByEmail(email string) (*entities.User, *e.Error) {
	//check if email is registered
	doEmailExist, err := uc.userRepo.IsEmailRegistered(email)
	if err != nil {
		return nil, err
	}
	if !doEmailExist {
		return nil, errEmailNotRegistered_401
	}

	//get user
	return uc.userRepo.GetUserByEmail(email)
}

func (uc *UserUseCase) SendOtpForPWChange(user *entities.User) (*string, *e.Error) {
	if err := otpManager.SendOtp(user.Phone); err != nil {
		return nil, e.SetError("Error while sending otp", err, 500)
	}
	resetPWClaims := struct {
		ID     uint
		Phone  string
		Status string
	}{
		ID:     user.ID,
		Phone:  user.Phone,
		Status: "PW change requested, otp not verified",
	}
	tokenString, errr := jwttoken.GenerateToken("user", resetPWClaims, time.Minute*5)
	if errr != nil {
		return nil, e.SetError("error generating token", errr, 500)
	}

	return &tokenString, nil
}
func (uc *UserUseCase) VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, *e.Error) {
	matchStatus, err := otpManager.VerifyOtp(phone, otp)
	if err != nil {
		return false, nil, e.SetError("Error while verifying otp", err, 500)
	}
	if !matchStatus {
		return false, nil, nil
	}
	resetPWClaims := struct {
		ID     uint
		Phone  string
		Status string
	}{
		ID:     id,
		Phone:  phone,
		Status: "PW change requested, otp verified",
	}
	tokenString, errr := jwttoken.GenerateToken("user", resetPWClaims, time.Minute*60)
	if errr != nil {
		return false, nil, e.SetError("error generating token", errr, 500)
	}
	return true, &tokenString, nil

}

func (uc *UserUseCase) ResetPasswordToNewPassword(id uint, newPassword *string) *e.Error {
	hashedPwd, err := hashpassword.Hashpassword(*newPassword)
	if err != nil {
		return e.SetError("error while hashing pw", err, 500)
	}
	return uc.userRepo.ResetPassword(id, &hashedPwd)
}

func (uc *UserUseCase) SetInitialPassword(id uint, newPassword *string) (string, *e.Error) {
	hashedPwd, errr:= hashpassword.Hashpassword(*newPassword)
	if errr != nil {
		return "", e.SetError("error while hashing pw", errr, 500)
	}

	err:= uc.userRepo.ResetPassword(id, &hashedPwd)
	if err != nil {
		return "", e.SetError("error while setting initial password", err, 500)
	}

	user, err := uc.userRepo.GetUserDetailsByID(id)
	if err != nil {
		return "", e.SetError("error while getting user details", err, 500)
	}

	newToken, errr := jwttoken.GenerateToken("user", user, time.Hour*24*30)
	if errr != nil {
		return "", e.SetError("error while generating token", err, 500)
	}

	return newToken, nil
}
