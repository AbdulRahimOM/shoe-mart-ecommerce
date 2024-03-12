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
	"fmt"
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
		return &e.Error{Err: errors.New(err.Error() + "Error occured while sending otp"), StatusCode: 500}
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
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- email is not registered\n.")
		return nil, &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 400}
	}

	//get userpassword from database
	userForToken, err := uc.userRepo.GetPasswordAndUserDetailsByEmail(req.Email)
	if err != nil {
		fmt.Println("Error occured while getting password from record")
		return nil, err
	}

	//check for password
	if hashpassword.CompareHashedPassword(userForToken.Password, req.Password) != nil {
		fmt.Println("Password Mismatch")
		return nil, &e.Error{Err: e.ErrInvalidPassword, StatusCode: 400}
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
	// fmt.Println("-----\nreq.email:", req.Email, "\n------")
	emailAlreadyUsed, err := uc.userRepo.IsEmailRegistered(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if emailAlreadyUsed {
		fmt.Println("\n email is already used!!\n.")
		return nil, &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 400}
	}

	hashedPwd, errr := hashpassword.Hashpassword(req.Password)
	if errr != nil {
		errr = errors.New("error while hashing pw. Error:)" + err.Error())
		return nil, &e.Error{Err: errr, StatusCode: 500}
	}
	fmt.Println("\n\n\n\nhashedpw=", hashedPwd)
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
		return &e.Error{Err: errors.New("address name already exists"), StatusCode: 400} //P-update
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
	//check if address name exists by ID
	doAddressExistsByID, err := uc.userRepo.DoAddressExistsByID(req.ID)
	if err != nil {
		return err
	}
	if !doAddressExistsByID {
		return &e.Error{Err: errors.New("address ID doesn't exist"), StatusCode: 400}
	}

	//check if userid is not changed
	userID, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		fmt.Println("couldn't get userID")
		return err
	}
	if userID != req.UserID {
		return &e.Error{Err: errors.New("UserID is not same, Corrupt request"), StatusCode: 401}
	}

	//get old address name
	oldAddressName, err := uc.userRepo.GetAddressNameByID(req.ID)
	if err != nil {
		fmt.Println("Error occured while getting old address name")
		return err
	}
	if oldAddressName != req.AddressName {
		//check if AddressNameExistsForAnotherAddressOfUser,
		doAddressNameAlreadyExists, err := uc.userRepo.DoAddressNameExists(req.AddressName)
		if err != nil {
			return err
		}
		if doAddressNameAlreadyExists {
			return &e.Error{Err: errors.New("address name already exists for user's another address"), StatusCode: 400}
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
	//check if address exists by ID
	doAddressExistsByID, err := uc.userRepo.DoAddressExistsByID(req.ID)
	if err != nil {
		fmt.Println("Error occured while checking if address exists")
		return err
	}
	if !doAddressExistsByID {
		return &e.Error{Err: errors.New("address ID doesn't exist"), StatusCode: 400}
	}

	//check if userid is conforming
	userID, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		fmt.Println("couldn't get userID")
		return err
	}
	if userID != req.UserID {
		return &e.Error{Err: errors.New("UserID is not same, Corrupt request"), StatusCode: 401}
	}

	//delete address
	return uc.userRepo.DeleteUserAddress(req.ID)
}

// GetUserAddresses
func (uc *UserUseCase) GetUserAddresses(userID uint) (*[]entities.UserAddress, *e.Error) {
	//check if user exists //this is required because if user is deleted from records, jwt token may not have expired.
	doUserExists, err := uc.userRepo.DoUserExistsByID(userID)
	if err != nil {
		return nil, err
	}
	if !doUserExists {
		return nil, &e.Error{Err: errors.New("user ID doesn't exist"), StatusCode: 400}
	}

	//get user addresses
	return uc.userRepo.GetUserAddresses(userID)
}

// GetProfile
func (uc *UserUseCase) GetProfile(userID uint) (*entities.UserDetails, *e.Error) {
	//check if user exists //this is required because if user is deleted from records, jwt token may not have expired.
	doUserExists, err := uc.userRepo.DoUserExistsByID(userID)
	if err != nil {
		return nil, err
	}
	if !doUserExists {
		fmt.Println("User ID doesn't exist")
		return nil, &e.Error{Err: errors.New("user ID doesn't exist"), StatusCode: 400}
	}

	//get user profile
	return uc.userRepo.GetProfile(userID)
}

// EditProfile
func (uc *UserUseCase) EditProfile(userID uint, req *request.EditProfileReq) *e.Error {
	//check if user exists //this is required because if user is deleted from records, jwt token may not have expired.
	doUserExists, err := uc.userRepo.DoUserExistsByID(userID)
	if err != nil {
		fmt.Println("Error occured while checking if user exists")
		return err
	}
	if !doUserExists {
		return &e.Error{Err: errors.New("user ID doesn't exist"), StatusCode: 400}
	}

	//check if email is not changed
	oldEmail, err := uc.userRepo.GetEmailByID(userID)
	if err != nil {
		fmt.Println("Error occured while getting old email")
		return err
	}
	if oldEmail != req.Email {
		//check if email already exists
		doEmailAlreadyExists, err := uc.userRepo.IsEmailRegistered(req.Email)
		if err != nil {
			fmt.Println("Error occured while checking if email already exists")
			return err
		}
		if doEmailAlreadyExists {
			return &e.Error{Err: errors.New("email already exists for another user"), StatusCode: 400}
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
		fmt.Println("Error occured while checking if email already exists")
		return nil, err
	}
	if !doEmailExist {
		fmt.Println("Email doesn't exist")
		return nil, &e.Error{Err: errors.New("email doesn't exist"), StatusCode: 400}
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
		fmt.Println("Error occured while verifying otp")
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
		fmt.Println("\n error while hashing pw. Error:)", err, "\n.")
		return &e.Error{Err: err, StatusCode: 500}
	}
	return uc.userRepo.ResetPassword(id, &hashedPwd)
}
