package accountsUsecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	hashpassword "MyShoo/pkg/hash_Password"
	jwttoken "MyShoo/pkg/jwt_tokens"
	otpManager "MyShoo/pkg/twilio"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
)

type UserUseCase struct {
	userRepo repoInterface.IUserRepo
}

func NewUserUseCase(repo repoInterface.IUserRepo) usecaseInterface.IUserUC {
	return &UserUseCase{userRepo: repo}
}

func (uc *UserUseCase) SendOtp(phone string) error {
	err := otpManager.SendOtp(phone)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) VerifyOtp(phone string, email string, otp string) (bool, error) {
	matchStatus, err := otpManager.VerifyOtp(phone, otp)
	if err != nil {
		fmt.Println("Error occured while verifying otp")
		return false, err
	}
	if !matchStatus {
		fmt.Println("otp mismatch")
		return false, errors.New("otp mismatch")
	} else {
		fmt.Println("otp matched")
		err := uc.userRepo.UpdateUserStatus(email, "verified")
		if err != nil {
			fmt.Println("Error occured while updating user status")
			return true, err
		}
		return true, nil
	}
}
func (uc *UserUseCase) SignIn(req *requestModels.UserSignInReq) (*string, error) {
	// fmt.Println("req.email=", req.Email)
	isEmailRegistered, err := uc.userRepo.IsEmailRegistered(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- email is not registered\n.")
		return nil, e.ErrEmailNotRegistered
	}

	//get userpassword from database
	hashedPassword, userInToken, err := uc.userRepo.GetPasswordAndUserDetailsByEmail(req.Email)
	if err != nil {
		fmt.Println("Error occured while getting password from record")
		return nil, err
	}

	//check for password
	if hashpassword.CompareHashedPassword(hashedPassword, req.Password) != nil {
		fmt.Println("Password Mismatch")
		return nil, e.ErrInvalidPassword
	}

	//generate token
	tokenString, err := jwttoken.GenerateToken("user", userInToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}

	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}

func (uc *UserUseCase) SignUp(req *requestModels.UserSignUpReq) (*string, error) {
	// fmt.Println("-----\nreq.email:", req.Email, "\n------")
	emailAlreadyUsed, err := uc.userRepo.IsEmailRegistered(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if emailAlreadyUsed {
		fmt.Println("\n email is already used!!\n.")
		return nil, e.ErrEmailAlreadyUsed
	}

	hashedPwd, err := hashpassword.Hashpassword(req.Password)
	if err != nil {
		fmt.Println("\n error while hashing pw. Error:)", err, "\n.")
		return nil, err
	}

	var signingUser = entities.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPwd,
		Status:    "not verified",
	}

	err = uc.userRepo.CreateUser(&signingUser)
	if err != nil {
		fmt.Println("\nUC: error recieved from ~repo.createuser()")
		return nil, err
	}

	var userInToken = entities.UserDetails{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "not verified",
	}

	//generate token
	tokenString, err := jwttoken.GenerateToken("user", userInToken, time.Hour*5)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (uc *UserUseCase) AddUserAddress(req *requestModels.AddUserAddress) error {
	//check if address name already exists
	doAddressNameAlreadyExists, err := uc.userRepo.DoAddressNameExists(req.AddressName)
	if err != nil {
		fmt.Println("Error occured while checking if address name exists")
		return err
	}
	if doAddressNameAlreadyExists {
		fmt.Println("An address already exists for user")
		return errors.New("address name already exists")
	}
	var address entities.UserAddress
	if err := copier.Copy(&address, &req); err != nil {
		fmt.Println("Error occured while copying request to cart entity")
		return err
	}
	//add address
	if err := uc.userRepo.AddUserAddress(&address); err != nil {
		return err
	}

	return nil
}

// EditUserAddress
func (uc *UserUseCase) EditUserAddress(req *requestModels.EditUserAddress) error {
	//check if address name exists by ID
	doAddressExistsByID, err := uc.userRepo.DoAddressExistsByID(req.ID)
	if err != nil {
		fmt.Println("Error occured while checking if address name exists")
		return err
	}
	if !doAddressExistsByID {
		fmt.Println("Address ID doesn't exist")
		return errors.New("address ID doesn't exist")
	}

	//check if userid is not changed
	userID, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		fmt.Println("couldn't get userID")
		return err
	}
	if userID != req.UserID {
		fmt.Println("UserID is not same, Corrupt request")
		return errors.New("UserID is not same, Corrupt request")
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
			fmt.Println("Error occured while checking if address name already exists for another address")
			return err
		}
		if doAddressNameAlreadyExists {
			fmt.Println("An address already exists for user's another address")
			return errors.New("address name already exists for user's another address")
		}
	}

	var address entities.UserAddress
	if err := copier.Copy(&address, &req); err != nil {
		fmt.Println("Error occured while copying request to cart entity")
		return err
	}
	//edit address
	if err := uc.userRepo.EditUserAddress(&address); err != nil {
		return err
	}

	return nil
}

// DeleteUserAddress
func (uc *UserUseCase) DeleteUserAddress(req *requestModels.DeleteUserAddress) error {
	//check if address exists by ID
	doAddressExistsByID, err := uc.userRepo.DoAddressExistsByID(req.ID)
	if err != nil {
		fmt.Println("Error occured while checking if address exists")
		return err
	}
	if !doAddressExistsByID {
		fmt.Println("Address ID doesn't exist")
		return errors.New("address ID doesn't exist")
	}

	//check if userid is conforming
	userID, err := uc.userRepo.GetUserIDFromAddressID(req.ID)
	if err != nil {
		fmt.Println("couldn't get userID")
		return err
	}
	if userID != req.UserID {
		fmt.Println("UserID is not same, Corrupt request")
		return errors.New("UserID is not same, Corrupt request")
	}

	//delete address
	if err := uc.userRepo.DeleteUserAddress(req.ID); err != nil {
		return err
	}

	return nil
}

// GetUserAddresses
func (uc *UserUseCase) GetUserAddresses(userID uint) (*[]entities.UserAddress, error) {
	//check if user exists //this is required because if user is deleted from records, jwt token may not have expired.
	doUserExists, err := uc.userRepo.DoUserExistsByID(userID)
	if err != nil {
		fmt.Println("Error occured while checking if user exists")
		return nil, err
	}
	if !doUserExists {
		fmt.Println("User ID doesn't exist")
		return nil, errors.New("user ID doesn't exist")
	}

	//get user addresses
	addresses, err := uc.userRepo.GetUserAddresses(userID)
	if err != nil {
		fmt.Println("Error occured while getting user addresses")
		return nil, err
	}

	return addresses, nil
}

// GetProfile
func (uc *UserUseCase) GetProfile(userID uint) (*entities.UserDetails, error) {
	//check if user exists //this is required because if user is deleted from records, jwt token may not have expired.
	doUserExists, err := uc.userRepo.DoUserExistsByID(userID)
	if err != nil {
		fmt.Println("Error occured while checking if user exists")
		return nil, err
	}
	if !doUserExists {
		fmt.Println("User ID doesn't exist")
		return nil, errors.New("user ID doesn't exist")
	}

	//get user profile
	var user *entities.UserDetails
	user, err = uc.userRepo.GetProfile(userID)
	if err != nil {
		fmt.Println("Error occured while getting user profile")
		return nil, err
	}

	return user, nil
}

// EditProfile
func (uc *UserUseCase) EditProfile(userID uint, req *requestModels.EditProfileReq) error {
	//check if user exists //this is required because if user is deleted from records, jwt token may not have expired.
	doUserExists, err := uc.userRepo.DoUserExistsByID(userID)
	if err != nil {
		fmt.Println("Error occured while checking if user exists")
		return err
	}
	if !doUserExists {
		fmt.Println("User ID doesn't exist")
		return errors.New("user ID doesn't exist")
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
			fmt.Println("Email already exists for another user")
			return errors.New("email already exists for another user")
		}
	}

	//edit profile
	if err := uc.userRepo.EditProfile(userID, req); err != nil {
		return err
	}

	return nil
}

// GetUserPhoneByEmail
func (uc *UserUseCase) GetUserByEmail(email string) (*entities.User, error) {
	//check if email is registered
	doEmailExist, err := uc.userRepo.IsEmailRegistered(email)
	if err != nil {
		fmt.Println("Error occured while checking if email already exists")
		return nil, err
	}
	if !doEmailExist {
		fmt.Println("Email doesn't exist")
		return nil, errors.New("email doesn't exist")
	}

	//get user
	var user *entities.User
	user, err = uc.userRepo.GetUserByEmail(email)
	if err != nil {
		fmt.Println("Error occured while getting user phone")
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) SendOtpForPWChange(user *entities.User) (*string, error) {
	if err := otpManager.SendOtp(user.Phone); err != nil {
		return nil, err
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
	tokenString, err := jwttoken.GenerateToken("user", resetPWClaims, time.Minute*5)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
func (uc *UserUseCase) VerifyOtpForPWChange(id uint, phone string, otp string) (bool, *string, error) {
	matchStatus, err := otpManager.VerifyOtp(phone, otp)
	if err != nil {
		fmt.Println("Error occured while verifying otp")
		return false, nil, err
	}
	if !matchStatus {
		fmt.Println("otp mismatch")
		return false, nil, errors.New("otp mismatch")
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
	tokenString, err := jwttoken.GenerateToken("user", resetPWClaims, time.Minute*60)
	if err != nil {
		return false, nil, err
	}
	return true, &tokenString, nil

}

func (uc *UserUseCase) ResetPassword(id uint, newPassword *string) error {
	hashedPwd, err := hashpassword.Hashpassword(*newPassword)
	if err != nil {
		fmt.Println("\n error while hashing pw. Error:)", err, "\n.")
		return err
	}
	if err := uc.userRepo.ResetPassword(id, &hashedPwd); err != nil {
		return err
	}
	return nil
}
