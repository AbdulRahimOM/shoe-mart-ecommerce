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
	var sellerlist *[]entities.PwMaskedSeller
	sellerlist, err := uc.adminRepo.GetSellersList()
	if err != nil {
		fmt.Println("Error occured while getting sellers list from database")
		return nil, err
	}
	return sellerlist, nil
}

func (uc *AdminUseCase) BlockUser(req *request.BlockUserReq) *e.Error {
	//check if user exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsUser(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- No such user (This email is not registered as a user)\n.")
		return e.ErrEmailNotRegistered
	}
	if err := uc.adminRepo.UpdateUserStatus(req.Email, "blocked"); err != nil {
		fmt.Println("Error occured while blocking user")
		return err
	}
	return nil
}
func (uc *AdminUseCase) UnblockUser(req *request.UnblockUserReq) *e.Error {
	//check if user exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsUser(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- No such user (This email is not registered as a user)\n.")
		return e.ErrEmailNotRegistered
	}

	if err := uc.adminRepo.UpdateUserStatus(req.Email, "verified"); err != nil {
		fmt.Println("Error occured while reverting blocked status to \"verified\" user")
		return err
	}
	return nil
}

func (uc *AdminUseCase) BlockSeller(req *request.BlockSellerReq) *e.Error {
	//check if seller exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsSeller(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return err
	}

	if !(isEmailRegistered) {
		fmt.Println("\n-- No such seller (This email is not registered as a seller)\n.")
		return e.ErrEmailNotRegistered
	}

	if err := uc.adminRepo.UpdateSellerStatus(req.Email, "blocked"); err != nil {
		fmt.Println("Error occured while blocking seller")
		return err
	}
	return nil
}

func (uc *AdminUseCase) UnblockSeller(req *request.UnblockSellerReq) *e.Error {
	//check if seller exists
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsSeller(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- No such seller (This email is not registered as a seller)\n.")
		return e.ErrEmailNotRegistered
	}

	if err := uc.adminRepo.UpdateSellerStatus(req.Email, "verified"); err != nil {
		fmt.Println("Error occured while reverting blocked status to \"verified\" seller")
		return err
	}
	return nil
}

func (uc *AdminUseCase) GetUsersList() (*[]entities.UserDetails, *e.Error) {
	userlist, err := uc.adminRepo.GetUsersList()
	if err != nil {
		fmt.Println("Error occured while getting users list from database")
		return nil, err
	}
	return userlist, nil
}
func (uc *AdminUseCase) SignIn(req *request.AdminSignInReq) (*string, *e.Error) {
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsAdmin(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- email is not registered\n.")
		return nil, &e.Error{Err:e.ErrEmailNotRegistered,StatusCode: 200}
	}

	//get adminpassword from database
	hashedPassword, adminInToken, err := uc.adminRepo.GetPasswordAndAdminDetailsByEmail(req.Email)
	if err != nil {
		fmt.Println("Error occured while getting password from record")
		return nil, err
	}

	//check for password
	if hashpassword.CompareHashedPassword(hashedPassword, req.Password) != nil {
		fmt.Println("Password Mismatch")
		return nil, &e.Error{Err:e.ErrInvalidPassword,StatusCode: 200}
	}

	//generate token
	tokenString, err2 := jwttoken.GenerateToken("admin", adminInToken, time.Hour*24*30)
	if err2 != nil {
		return nil, &e.Error{Err:err,StatusCode: 500}
	}
	return &tokenString, nil
}

func (uc *AdminUseCase) RestartConfig() error {	//err pro
	err := config.RestartDeliveryConfig()
	if err != nil {
		fmt.Println("Error occured while reloading config")
		return  err
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
		return errors.New("seller is already verified")
	}

	err = uc.adminRepo.VerifySeller(req.SellerID)
	if err != nil {
		fmt.Println("Error occured while verifying seller")
		return err
	}

	return nil
}
