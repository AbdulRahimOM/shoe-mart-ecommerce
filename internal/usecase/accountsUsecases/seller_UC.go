package accountsusecase

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
	hashpassword "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/hashPassword"
	jwttoken "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/jwt"

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
	if err := hashpassword.CompareHashedPassword(sellerForToken.Password, req.Password); err != nil {
		return nil, e.SetError("wrong password", err, 401)
	}

	//generate token
	tokenString, errr := jwttoken.GenerateToken("seller", sellerForToken, time.Hour*24*30)
	if errr != nil {
		return nil, e.SetError("error generating token", errr, 500)
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
		return nil, e.SetError("error generating token", errr, 500)
	}
	// fmt.Println("token created + 'BNo error' is sending to handler layer")
	return &tokenString, nil
}
