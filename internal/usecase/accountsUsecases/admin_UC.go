package accountsusecase

import (
	"MyShoo/internal/config"
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	hashpassword "MyShoo/pkg/hashPassword"
	jwttoken "MyShoo/pkg/jwt"
	"errors"
	"fmt"
	"time"
)

type AdminUseCase struct {
	adminRepo repoInterface.IAdminRepo
}

func NewAdminUseCase(repo repoInterface.IAdminRepo) usecase.IAdminUC {
	return &AdminUseCase{adminRepo: repo}
}

func (uc *AdminUseCase) GetSellersList() (*[]entities.PwMaskedSeller, *e.Error) {
	return uc.adminRepo.GetSellersList()
}

func (uc *AdminUseCase) BlockUser(req *request.BlockUserReq) *e.Error {
	//check if user exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsUser(req.Email)
	if err != nil {
		return err
	}
	if !(isEmailRegistered) {
		return &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 400}
	}
	return uc.adminRepo.UpdateUserStatus(req.Email, "blocked")
}

func (uc *AdminUseCase) UnblockUser(req *request.UnblockUserReq) *e.Error {
	//check if user exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsUser(req.Email)
	if err != nil {
		return err
	}
	if !(isEmailRegistered) {
		return &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 400}
	}
	return uc.adminRepo.UpdateUserStatus(req.Email, "verified")
}

func (uc *AdminUseCase) BlockSeller(req *request.BlockSellerReq) *e.Error {
	//check if seller exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsSeller(req.Email)
	if err != nil {
		return err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- No such seller (This email is not registered as a seller)\n.")
		return &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 400}
	}
	return uc.adminRepo.UpdateSellerStatus(req.Email, "blocked")
}

func (uc *AdminUseCase) UnblockSeller(req *request.UnblockSellerReq) *e.Error {
	//check if seller exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsSeller(req.Email)
	if err != nil {
		return err
	}
	if !(isEmailRegistered) {
		return &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 400}
	}

	return uc.adminRepo.UpdateSellerStatus(req.Email, "verified")
}

func (uc *AdminUseCase) GetUsersList() (*[]entities.UserDetails, *e.Error) {
	return uc.adminRepo.GetUsersList()
}

func (uc *AdminUseCase) SignIn(req *request.AdminSignInReq) (*string, *e.Error) {
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsAdmin(req.Email)
	if err != nil {
		return nil, err
	}
	if !(isEmailRegistered) {
		return nil, &e.Error{Err: e.ErrEmailNotRegistered, StatusCode: 200}
	}

	//get adminpassword from database
	hashedPassword, adminInToken, err := uc.adminRepo.GetPasswordAndAdminDetailsByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//check for password
	if hashpassword.CompareHashedPassword(hashedPassword, req.Password) != nil {
		return nil, &e.Error{Err: e.ErrInvalidPassword, StatusCode: 200}
	}

	//generate token
	tokenString, err2 := jwttoken.GenerateToken("admin", adminInToken, time.Hour*24*30)
	if err2 != nil {
		return nil, &e.Error{Err: err, StatusCode: 500}
	}
	return &tokenString, nil
}

func (uc *AdminUseCase) RestartConfig() *e.Error { //err pro
	err := config.RestartDeliveryConfig()
	if err != nil {
		fmt.Println("Error occured while reloading config")
		return &e.Error{Err: err, StatusCode: 500}
	}
	return nil
}

// VerifySeller
func (uc *AdminUseCase) VerifySeller(req *request.VerifySellerReq) *e.Error {
	//check if status is not verified
	isVerified, err := uc.adminRepo.IsSellerVerified(req.SellerID)
	if err != nil {
		fmt.Println("Error occured while getting seller status")
		return err
	}
	if isVerified {
		fmt.Println("\n-- seller is already verified\n.")
		return &e.Error{Err: errors.New("seller is already verified"), StatusCode: 400}
	}

	return uc.adminRepo.VerifySeller(req.SellerID)
}
