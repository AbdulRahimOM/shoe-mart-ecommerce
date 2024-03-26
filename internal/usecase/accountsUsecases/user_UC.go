package accountsusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	hashpassword "MyShoo/pkg/hashPassword"
	jwttoken "MyShoo/pkg/jwt"
	otpManager "MyShoo/pkg/twilio"
	"errors"
	"time"

	"github.com/jinzhu/copier"
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
		return e.SetError("Error occured while sending otp", err, 500)
	}
	return nil
}

func (uc *UserUseCase) VerifyOtp(phone string, email string, otp string) (bool, *e.Error) {
	if matchStatus, err := otpManager.VerifyOtp(phone, otp); err != nil {
		return false, &e.Error{Err: err, StatusCode: 500}
	} else if !matchStatus {
		return false, nil
	} else {
		//otp matched!!!, update status
		return true, uc.userRepo.UpdateUserStatus(email, "verified")
	}
}
func (uc *UserUseCase) SignIn(req *request.UserSignInReq) (*string, *e.Error) {
	// fmt.Println("req.email=", req.Email)
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
	if hashpassword.CompareHashedPassword(userForToken.Password, req.Password) != nil {
		return nil, e.ErrInvalidPassword_401
	}

	//generate token
	tokenString, errr := jwttoken.GenerateToken("user", userForToken, time.Hour*24*30)
	if errr != nil {
		return nil, &e.Error{Err: errr, StatusCode: 500}
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
		errr = errors.New("error while hashing pw. Error:)" + err.Error())
		return nil, &e.Error{Err: errr, StatusCode: 500}
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
		return nil, &e.Error{Err: errr, StatusCode: 500}
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
		return e.SetError("address name already exists",nil, 400) 
	}
	var address entities.UserAddress
	if err := copier.Copy(&address, &req); err != nil {
		return &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to cart entity"), StatusCode: 500}
	}
	//add address
	return uc.userRepo.AddUserAddress(&address)
}

// EditUserAddress
func (uc *UserUseCase) EditUserAddress(req *request.EditUserAddress) *e.Error {
	//check if userid is conforming
	userID, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		return err
	}
	if userID != req.UserID {
		return e.SetError("UserID is not same, Corrupt request",nil, 401)
	}

	//get old address name
	oldAddressName, err := uc.userRepo.GetAddressNameByID(req.ID)
	if err != nil {
		return err
	}
	if oldAddressName != req.AddressName {
		//check if AddressNameExistsForAnotherAddressOfUser,
		doAddressNameAlreadyExists, err := uc.userRepo.DoAddressNameExists(req.AddressName)
		if err != nil {
			return err
		}
		if doAddressNameAlreadyExists {
			return e.SetError("address name already exists for user's another address",nil, 400)
		}
	}

	var address entities.UserAddress
	if err := copier.Copy(&address, &req); err != nil {
		return &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to cart entity"), StatusCode: 500}
	}
	//edit address
	return uc.userRepo.EditUserAddress(&address)
}

// DeleteUserAddress
func (uc *UserUseCase) DeleteUserAddress(req *request.DeleteUserAddress) *e.Error {
	//check if userid is conforming
	userID, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		return err
	}
	if userID != req.UserID {
		return e.SetError("UserID is not same, Corrupt request",nil, 401)
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
	if oldEmail != req.Email {
		//check if email already exists
		doEmailAlreadyExists, err := uc.userRepo.IsEmailRegistered(req.Email)
		if err != nil {
			return err
		}
		if doEmailAlreadyExists {
			return e.SetError("email already exists for another user",nil, 400)
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
		return nil, &e.Error{Err: err, StatusCode: 500}
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
		return nil, &e.Error{Err: errr, StatusCode: 500}
	}

	return &tokenString, nil
}
func (uc *UserUseCase) VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, *e.Error) {
	matchStatus, err := otpManager.VerifyOtp(phone, otp)
	if err != nil {
		return false, nil, &e.Error{Err: err, StatusCode: 500}
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
		return false, nil, &e.Error{Err: errr, StatusCode: 500}
	}
	return true, &tokenString, nil

}

func (uc *UserUseCase) ResetPassword(id uint, newPassword *string) *e.Error {
	hashedPwd, err := hashpassword.Hashpassword(*newPassword)
	if err != nil {
		return &e.Error{Err: err, StatusCode: 500}
	}
	return uc.userRepo.ResetPassword(id, &hashedPwd)
}
