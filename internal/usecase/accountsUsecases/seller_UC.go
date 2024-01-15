package accountsUsecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	hashpassword "MyShoo/pkg/hash_Password"
	jwttoken "MyShoo/pkg/jwt_tokens"

	"fmt"
	"time"
)

type SellerUseCase struct {
	sellerRepo repoInterface.ISellerRepo
}

func NewSellerUseCase(repo repoInterface.ISellerRepo) usecaseInterface.ISellerUC {
	return &SellerUseCase{sellerRepo: repo}
}

func (uc *SellerUseCase) SignIn(req *requestModels.SellerSignInReq) (*string, error) {
	// fmt.Println("req.email=", req.Email)
	isEmailRegistered, err := uc.sellerRepo.IsEmailRegistered(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- email is not registered\n.")
		return nil, e.ErrEmailNotRegistered
	}

	//get sellerpassword from database
	hashedPassword, sellerInToken, err := uc.sellerRepo.GetPasswordAndSellerDetailsByEmail(req.Email)
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
	tokenString, err := jwttoken.GenerateToken("seller", sellerInToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}

func (uc *SellerUseCase) SignUp(req *requestModels.SellerSignUpReq) (*string, error) {
	// fmt.Println("-----\nreq.email:", req.Email, "\n------")
	emailAlreadyUsed, err := uc.sellerRepo.IsEmailRegistered(req.Email)
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

	var signingSeller = entities.Seller{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPwd,
		Status:    "not verified",
	}

	err = uc.sellerRepo.CreateSeller(&signingSeller)
	if err != nil {
		fmt.Println("\nUC: error recieved from ~repo.createseller()")
		return nil, err
	}

	// //send OTP
	// err = otphelper.SendOtp(seller.Phone)
	// if err != nil {
	// 	uc.sellerRepo.DeleteByPhone(seller.Phone)
	// 	return nil, err
	// }
	var sellerInToken = entities.SellerDetails{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "not verified",
	}

	//generate token
	tokenString, err := jwttoken.GenerateToken("seller", sellerInToken, time.Hour*5)
	if err != nil {
		return nil, err
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}
