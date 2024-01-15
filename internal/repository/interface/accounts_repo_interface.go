package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
)

type IAdminRepo interface {
	IsEmailRegisteredAsAdmin(email string) (bool, error)
	GetPasswordAndAdminDetailsByEmail(email string) (string, entities.AdminDetails, error)
	GetUsersList() (*[]entities.UserDetails, error)
	UpdateUserStatus(email string, newStatus string) error
	GetSellersList() (*[]entities.SellerDetails, error)
	IsEmailRegisteredAsUser(email string) (bool, error)
	IsEmailRegisteredAsSeller(email string) (bool, error)
	UpdateSellerStatus(email string, newStatus string) error
}

type ISellerRepo interface {
	// SaveSellerData() error //inactive
	CreateSeller(*entities.Seller) error
	IsEmailRegistered(string) (bool, error)

	//returns hashed password, seller details as entities.SellerDetails struct and error (if any, else nil)
	GetPasswordAndSellerDetailsByEmail(string) (string, entities.SellerDetails, error)
}

type IUserRepo interface {
	// SaveUserData() error //inactive
	CreateUser(*entities.User) error
	IsEmailRegistered(string) (bool, error)
	GetUserByEmail(email string) (*entities.User, error)
	ResetPassword(id uint, newPassword *string) error

	//returns hashed password, user details as entities.UserDetails struct and error (if any, else nil)
	GetPasswordAndUserDetailsByEmail(string) (string, entities.UserDetails, error)
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
	GetProfile(userID uint) (*entities.UserDetails, error)
	GetEmailByID(userID uint) (string, error)
	EditProfile(userID uint, req *requestModels.EditProfileReq) error
}
