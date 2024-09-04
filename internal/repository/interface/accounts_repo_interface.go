package repo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
)

type IAdminRepo interface {
	IsEmailRegisteredAsAdmin(email string) (bool, *e.Error)
	GetPasswordAndAdminDetailsByEmail(email string) (*string, *entities.AdminDetails, *e.Error)
	GetUsersList() (*[]entities.UserDetails, *e.Error)
	UpdateUserStatus(userID uint, newStatus string) *e.Error
	GetSellersList() (*[]entities.PwMaskedSeller, *e.Error)
	UpdateSellerStatus(sellerID uint, newStatus string) *e.Error
	IsSellerVerified(sellerID uint) (bool, *e.Error)
	VerifySeller(sellerID uint) *e.Error
}

type ISellerRepo interface {
	// SaveSellerData() *e.Error //inactive
	CreateSeller(*entities.Seller) *e.Error
	IsEmailRegistered(string) (bool, *e.Error)

	//returns hashed password, seller details as entities.PwMaskedSeller struct and *e.Error (if any, else nil)
	GetSellerWithPwByEmail(string) (*entities.Seller, *e.Error)
}

// type IUserRepo interface {
// 	// SaveUserData() *e.Error //inactive
// 	CreateUser(*entities.User) *e.Error
// 	IsEmailRegistered(string) (bool, *e.Error)
// 	GetUserByEmail(email string) (*entities.User, *e.Error)
// 	GetUserBasicInfoByID(id uint) (*response.UserInfoForInvoice, *e.Error)
// 	ResetPassword(id uint, newPassword *string) *e.Error

// 	//returns hashed password, user details as entities.UserDetails struct and *e.Error (if any, else nil)
// 	GetPasswordAndUserDetailsByEmail(string) (*entities.User, *e.Error)
// 	UpdateUserStatus(email string, newStatus string) *e.Error

// 	AddUserAddress(newAddress *entities.UserAddress) *e.Error
// 	DoAddressNameExists(name string) (bool, *e.Error)
// 	EditUserAddress(newaddress *entities.UserAddress) *e.Error
// 	GetAddressNameByID(id uint) (string, *e.Error)
// 	DeleteUserAddress(id uint) *e.Error
// 	GetUserIDFromAddressID(id uint) (uint, *e.Error)

// 	GetUserAddresses(userId uint) (*[]entities.UserAddress, *e.Error)
// 	GetUserAddress(addressID uint) (*entities.UserAddress, *e.Error)

// 	GetProfile(userID uint) (*entities.UserDetails, *e.Error)
// 	GetEmailByUserID(userID uint) (string, *e.Error)
// 	EditProfile(userID uint, req *request.EditProfileReq) *e.Error

// 	GetWalletBalance(userID uint) (float32, *e.Error)
// }
