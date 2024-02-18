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

func (uc *AdminUseCase) GetSellersList() (*[]entities.PwMaskedSeller, error) {
	var sellerlist *[]entities.PwMaskedSeller
	sellerlist, err := uc.adminRepo.GetSellersList()
	if err != nil {
		fmt.Println("Error occured while getting sellers list from database")
		return nil, err
	}
	return sellerlist, nil
}

func (uc *AdminUseCase) BlockUser(req *request.BlockUserReq) error {
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
func (uc *AdminUseCase) UnblockUser(req *request.UnblockUserReq) error {
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

func (uc *AdminUseCase) BlockSeller(req *request.BlockSellerReq) error {
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

func (uc *AdminUseCase) UnblockSeller(req *request.UnblockSellerReq) error {
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

func (uc *AdminUseCase) GetUsersList() (*[]entities.UserDetails, error) {
	userlist, err := uc.adminRepo.GetUsersList()
	if err != nil {
		fmt.Println("Error occured while getting users list from database")
		return nil, err
	}
	return userlist, nil
}
func (uc *AdminUseCase) SignIn(req *request.AdminSignInReq) (*string, error) {
	isEmailRegistered, err := uc.adminRepo.IsEmailRegisteredAsAdmin(req.Email)
	if err != nil {
		fmt.Println("Error occured while searching email, error:", err)
		return nil, err
	}
	if !(isEmailRegistered) {
		fmt.Println("\n-- email is not registered\n.")
		return nil, e.ErrEmailNotRegistered
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
		return nil, e.ErrInvalidPassword
	}

	//generate token
	tokenString, err := jwttoken.GenerateToken("admin", adminInToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (uc *AdminUseCase) RestartConfig() error {
	err := config.RestartDeliveryConfig()
	if err != nil {
		fmt.Println("Error occured while reloading config")
		return err
	}
	return nil
}

// VerifySeller
func (uc *AdminUseCase) VerifySeller(req *request.VerifySellerReq) error {
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
