package accountsusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	hashpassword "MyShoo/pkg/hashPassword"
	jwttoken "MyShoo/pkg/jwt"

	"fmt"
	"time"
)

type SellerUseCase struct {
	sellerRepo repoInterface.ISellerRepo
}

func NewSellerUseCase(repo repoInterface.ISellerRepo) usecase.ISellerUC {
	return &SellerUseCase{sellerRepo: repo}
}

func (uc *SellerUseCase) SignIn(req *request.SellerSignInReq) (*string, error) {
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
	sellerForToken, err := uc.sellerRepo.GetSellerWithPwByEmail(req.Email)
	if err != nil {
		fmt.Println("Error occured while getting password from record")
		return nil, err
	}

	//check for password
	if hashpassword.CompareHashedPassword(sellerForToken.Password, req.Password) != nil {
		fmt.Println("Password Mismatch")
		return nil, e.ErrInvalidPassword
	}

	//generate token
	tokenString, err := jwttoken.GenerateToken("seller", sellerForToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}

func (uc *SellerUseCase) SignUp(req *request.SellerSignUpReq) (*string, error) {
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
	var sellerForToken = entities.PwMaskedSeller{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "not verified",
	}

	//generate token
	tokenString, err := jwttoken.GenerateToken("seller", sellerForToken, time.Hour*5)
	if err != nil {
		return nil, err
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}
