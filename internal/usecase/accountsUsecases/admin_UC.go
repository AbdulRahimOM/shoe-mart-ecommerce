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
	"time"
)

var (
	errEmailNotRegistered_401      = e.ErrEmailNotRegistered_401
	errEmailAlreadyUsed_401        = e.ErrEmailAlreadyUsed_401
	errSellerIsAlreadyVerified_400 = &e.Error{Status: "failed",Msg: "Seller is already verified", Err: nil, StatusCode: 400}
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
	return uc.adminRepo.UpdateUserStatus(req.UserID, "blocked")
}

func (uc *AdminUseCase) UnblockUser(req *request.UnblockUserReq) *e.Error {
	return uc.adminRepo.UpdateUserStatus(req.UserID, "verified")
}

func (uc *AdminUseCase) BlockSeller(req *request.BlockSellerReq) *e.Error {
	return uc.adminRepo.UpdateSellerStatus(req.SellerID, "blocked")
}

func (uc *AdminUseCase) UnblockSeller(req *request.UnblockSellerReq) *e.Error {
	return uc.adminRepo.UpdateSellerStatus(req.SellerID, "verified")
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
		return nil, errEmailNotRegistered_401
	}

	//get adminpassword from database
	hashedPassword, adminInToken, err := uc.adminRepo.GetPasswordAndAdminDetailsByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//check for password
	if err:=hashpassword.CompareHashedPassword(*hashedPassword, req.Password);err != nil {
		return nil, e.SetError("password mismatch", err, 401)
	}

	//generate token
	tokenString, err2 := jwttoken.GenerateToken("admin", adminInToken, time.Hour*24*30)
	if err2 != nil {
		return nil, e.SetError("error generating token", err2, 500)
	}
	return &tokenString, nil
}

func (uc *AdminUseCase) RestartConfig() *e.Error { //err pro
	err := config.RestartDeliveryConfig()
	if err != nil {
		return e.SetError("error restarting config", err, 500)
	}
	return nil
}

// VerifySeller
func (uc *AdminUseCase) VerifySeller(req *request.VerifySellerReq) *e.Error {
	//check if status is not verified
	isVerified, err := uc.adminRepo.IsSellerVerified(req.SellerID)
	if err != nil {
		return err
	}
	if isVerified {
		return errSellerIsAlreadyVerified_400
	}

	return uc.adminRepo.VerifySeller(req.SellerID)
}
