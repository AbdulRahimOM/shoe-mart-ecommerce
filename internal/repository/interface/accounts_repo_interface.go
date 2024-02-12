package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type IAdminRepo interface {
	IsEmailRegisteredAsAdmin(email string) (bool, error)
	GetPasswordAndAdminDetailsByEmail(email string) (string, entities.AdminDetails, error)
	GetUsersList() (*[]entities.UserDetails, error)
	UpdateUserStatus(email string, newStatus string) error
	GetSellersList() (*[]entities.PwMaskedSeller, error)
	IsEmailRegisteredAsUser(email string) (bool, error)
	IsEmailRegisteredAsSeller(email string) (bool, error)
	UpdateSellerStatus(email string, newStatus string) error
	IsSellerVerified(sellerID uint) (bool, error)
	VerifySeller(sellerID  uint) error
}

type ISellerRepo interface {
	// SaveSellerData() error //inactive
	CreateSeller(*entities.Seller) error
	IsEmailRegistered(string) (bool, error)

	//returns hashed password, seller details as entities.PwMaskedSeller struct and error (if any, else nil)
	GetSellerWithPwByEmail(string) (*entities.Seller, error)
}

type IUserRepo interface {
	// SaveUserData() error //inactive
	CreateUser(*entities.User) error
	IsEmailRegistered(string) (bool, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserBasicInfoByID(id uint) (*response.UserInfoForInvoice, error)
	ResetPassword(id uint, newPassword *string) error

	//returns hashed password, user details as entities.UserDetails struct and error (if any, else nil)
	GetPasswordAndUserDetailsByEmail(string) (*entities.User, error)
	UpdateUserStatus(email string, newStatus string) error

	AddUserAddress(newAddress *entities.UserAddress) error
	DoAddressNameExists(name string) (bool, error)
	DoAddressExistsByID(id uint) (bool, error)
	DoAddressExistsByIDForUser(id uint, userID uint) (bool, error)
	EditUserAddress(newaddress *entities.UserAddress) error
	GetAddressNameByID(id uint) (string, error)
	DeleteUserAddress(id uint) error
	GetUserIDFromAddressID(id uint) (uint, error)

	DoUserExistsByID(userId uint) (bool, error)
	GetUserAddresses(userId uint) (*[]entities.UserAddress, error)
	GetUserAddress(addressID uint) (*entities.UserAddress, error)

	GetProfile(userID uint) (*entities.UserDetails, error)
	GetEmailByID(userID uint) (string, error)
	EditProfile(userID uint, req *requestModels.EditProfileReq) error

	GetWalletBalance(userID uint) (float32, error)
}
