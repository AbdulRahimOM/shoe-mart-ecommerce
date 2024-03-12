package repo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type IAdminRepo interface {
	IsEmailRegisteredAsAdmin(email string) (bool, *e.Error)
	GetPasswordAndAdminDetailsByEmail(email string) (string, entities.AdminDetails, *e.Error)
	GetUsersList() (*[]entities.UserDetails, *e.Error)
	UpdateUserStatus(email string, newStatus string) error
	GetSellersList() (*[]entities.PwMaskedSeller, *e.Error)
	IsEmailRegisteredAsUser(email string) (bool, *e.Error)
	IsEmailRegisteredAsSeller(email string) (bool, *e.Error)
	UpdateSellerStatus(email string, newStatus string) error
	IsSellerVerified(sellerID uint) (bool, *e.Error)
	VerifySeller(sellerID uint) error
}

type ISellerRepo interface {
	// SaveSellerData() error //inactive
	CreateSeller(*entities.Seller) error
	IsEmailRegistered(string) (bool, *e.Error)

	//returns hashed password, seller details as entities.PwMaskedSeller struct and error (if any, else nil)
	GetSellerWithPwByEmail(string) (*entities.Seller, *e.Error)
}

type IUserRepo interface {
	// SaveUserData() error //inactive
	CreateUser(*entities.User) error
	IsEmailRegistered(string) (bool, *e.Error)
	GetUserByEmail(email string) (*entities.User, *e.Error)
	GetUserBasicInfoByID(id uint) (*response.UserInfoForInvoice, *e.Error)
	ResetPassword(id uint, newPassword *string) error

	//returns hashed password, user details as entities.UserDetails struct and error (if any, else nil)
	GetPasswordAndUserDetailsByEmail(string) (*entities.User, *e.Error)
	UpdateUserStatus(email string, newStatus string) error

	AddUserAddress(newAddress *entities.UserAddress) error
	DoAddressNameExists(name string) (bool, *e.Error)
	DoAddressExistsByID(id uint) (bool, *e.Error)
	DoAddressExistsByIDForUser(id uint, userID uint) (bool, *e.Error)
	EditUserAddress(newaddress *entities.UserAddress) error
	GetAddressNameByID(id uint) (string, *e.Error)
	DeleteUserAddress(id uint) error
	GetUserIDFromAddressID(id uint) (uint, *e.Error)

	DoUserExistsByID(userId uint) (bool, *e.Error)
	GetUserAddresses(userId uint) (*[]entities.UserAddress, *e.Error)
	GetUserAddress(addressID uint) (*entities.UserAddress, *e.Error)

	GetProfile(userID uint) (*entities.UserDetails, *e.Error)
	GetEmailByID(userID uint) (string, *e.Error)
	EditProfile(userID uint, req *request.EditProfileReq) error

	GetWalletBalance(userID uint) (float32, *e.Error)
}
