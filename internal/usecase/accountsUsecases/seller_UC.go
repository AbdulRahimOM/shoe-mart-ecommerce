package accountsusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	hashpassword "MyShoo/pkg/hashPassword"
	jwttoken "MyShoo/pkg/jwt"

	"time"
)

type SellerUseCase struct {
	sellerRepo repoInterface.ISellerRepo
}

func NewSellerUseCase(repo repoInterface.ISellerRepo) usecase.ISellerUC {
	return &SellerUseCase{sellerRepo: repo}
}

func (uc *SellerUseCase) SignIn(req *request.SellerSignInReq) (*string, *e.Error) {
	// fmt.Println("req.email=", req.Email)
	isEmailRegistered, err := uc.sellerRepo.IsEmailRegistered(req.Email)
	if err != nil {
		return nil, err
	}
	if !(isEmailRegistered) {
		return nil, errEmailNotRegistered_401
	}

	//get sellerpassword from database
	sellerForToken, err := uc.sellerRepo.GetSellerWithPwByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//check for password
	if hashpassword.CompareHashedPassword(sellerForToken.Password, req.Password) != nil {
		return nil, &e.Error{Err: e.ErrInvalidPassword_401, StatusCode: 400}
	}

	//generate token
	tokenString, errr := jwttoken.GenerateToken("seller", sellerForToken, time.Hour*24*30)
	if errr != nil {
		return nil, &e.Error{Err: errr, StatusCode: 500}
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}

func (uc *SellerUseCase) SignUp(req *request.SellerSignUpReq) (*string, *e.Error) {
	// fmt.Println("-----\nreq.email:", req.Email, "\n------")
	emailAlreadyUsed, err := uc.sellerRepo.IsEmailRegistered(req.Email)
	if err != nil {
		return nil, err
	}
	if emailAlreadyUsed {
		return nil, errEmailNotRegistered_401
	}

	hashedPwd, errr := hashpassword.Hashpassword(req.Password)
	if errr != nil {
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
	tokenString, errr := jwttoken.GenerateToken("seller", sellerForToken, time.Hour*5)
	if errr != nil {
		return nil, &e.Error{Err: errr, StatusCode: 500}
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}
