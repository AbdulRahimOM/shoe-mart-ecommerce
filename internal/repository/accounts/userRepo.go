package accountrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repoInterface "MyShoo/internal/repository/interface"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) repoInterface.IUserRepo {
	return &UserRepo{DB: db}
}

func (repo *UserRepo) GetAddressNameByID(id uint) (string, *e.Error) {
	var name string
	query := repo.DB.Raw(`
		SELECT "addressName"
		FROM user_addresses
		WHERE id = ?`,
		id).Scan(&name)

	if query.Error != nil {
		return "", e.DBQueryError_500(&query.Error)
	}

	return name, nil
}

func (repo *UserRepo) EditUserAddress(newAddress *entities.UserAddress) *e.Error {
	result := repo.DB.Model(&entities.UserAddress{}).Where("id = ?", newAddress.ID).Updates(newAddress)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

func (repo *UserRepo) AddUserAddress(newAddress *entities.UserAddress) *e.Error {
	result := repo.DB.Create(&newAddress)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

func (repo *UserRepo) DoAddressNameExists(name string) (bool, *e.Error) {
	var temp entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE "addressName" = ?`,
		name).Scan(&temp)

	if query.Error != nil {
		return false, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *UserRepo) UpdateUserStatus(email string, newStatus string) *e.Error {
	var user entities.User
	err := repo.DB.Model(&user).Where("email = ?", email).Update("status", newStatus).Error
	if err != nil {
		return &e.Error{Err: err, StatusCode: 500}
	}

	return nil
}

func (repo *UserRepo) GetPasswordAndUserDetailsByEmail(email string) (*entities.User, *e.Error) {

	var user entities.User
	query := repo.DB.Raw(`
	SELECT * 
	FROM users 
	WHERE email = ?`,
		email).Scan(&user)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error) //existence
	}

	return &user, nil
}

func (repo *UserRepo) IsEmailRegistered(email string) (bool, *e.Error) {
	var emptyStruct struct{}
	query := repo.DB.Raw(`
        SELECT * 
        FROM users 
        WHERE email = ?`,
		email).Scan(&emptyStruct)

	if query.Error != nil {
		return false, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *UserRepo) CreateUser(user *entities.User) *e.Error {
	userCreation := repo.DB.Create(&user)
	if userCreation.Error != nil {
		return &e.Error{Err: userCreation.Error, StatusCode: 500}
	}
	return nil
}

// DeleteUserAddress
func (repo *UserRepo) DeleteUserAddress(id uint) *e.Error {
	result := repo.DB.Delete(&entities.UserAddress{}, id)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

// GetUserIDFromAddressID
func (repo *UserRepo) GetUserIDFromAddressID(id uint) (uint, *e.Error) {
	var userID uint
	query := repo.DB.Raw(`
		SELECT "userId"
		FROM user_addresses
		WHERE id = ?`,
		id).Scan(&userID)

	if query.Error != nil {
		return 0, e.DBQueryError_500(&query.Error)
	}
	if query.RowsAffected == 0 {
		return 0, e.SetError("no address found with this id", nil, 400)
	}

	return userID, nil
}

func (repo *UserRepo) GetUserAddresses(userId uint) (*[]entities.UserAddress, *e.Error) {
	var addresses []entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE "userId" = ?`,
		userId).Scan(&addresses)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}
	return &addresses, nil
}

func (repo *UserRepo) GetProfile(userID uint) (*entities.UserDetails, *e.Error) {
	var userDetails *entities.UserDetails
	query := repo.DB.Raw(`
		SELECT 
			id,
			"firstName",
			"lastName",
			email,
			phone,
			status
		FROM users
		WHERE id = ?`,
		userID).Scan(&userDetails)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, e.SetError("no user found with this id", nil,  400)
	}

	return userDetails, nil
}

// EditProfile implements repository_interface.IUserRepo.
func (repo *UserRepo) EditProfile(userID uint, req *request.EditProfileReq) *e.Error {

	result := repo.DB.Model(&entities.User{}).Where("id = ?", userID).Updates(req)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

// GetEmailByUserID implements repository_interface.IUserRepo.
func (repo *UserRepo) GetEmailByUserID(userID uint) (string, *e.Error) {
	var email string
	query := repo.DB.Raw(`
		SELECT email
		FROM users
		WHERE id = ?`,
		userID).Scan(&email)

	if query.Error != nil {
		return "", e.DBQueryError_500(&query.Error)
	}
	if query.RowsAffected == 0 {
		return "", e.SetError("no user found with this id", nil, 400)
	}

	return email, nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*entities.User, *e.Error) {
	var user entities.User
	query := repo.DB.Raw(`
		SELECT *
		FROM users
		WHERE email = ?`,
		email).Scan(&user)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &user, nil
}

func (repo *UserRepo) ResetPassword(id uint, newPassword *string) *e.Error {
	result := repo.DB.Model(&entities.User{}).Where("id = ?", id).Update("password", newPassword)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

// GetUserAddress(userID uint, addressID uint) (*entities.UserAddress, error)
func (repo *UserRepo) GetUserAddress(addressID uint) (*entities.UserAddress, *e.Error) {
	var address entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE id = ?`,
		addressID).Scan(&address)
	if query.Error != nil {
		return nil, &e.Error{Err: query.Error, StatusCode: 500}
	}
	if query.RowsAffected == 0 {
		return nil, e.SetError("no address found with this id", nil, 400)
	}

	return &address, nil
}

// GetWalletBalance
func (repo *UserRepo) GetWalletBalance(userID uint) (float32, *e.Error) {
	var balance float32
	query := repo.DB.Raw(`
		SELECT wallet_balance
		FROM users
		WHERE id = ?`,
		userID).Scan(&balance)
	if query.Error != nil {
		return 0, e.DBQueryError_500(&query.Error)
	}

	return balance, nil
}

// GetUserByID implements repository_interface.IUserRepo.
func (repo *UserRepo) GetUserBasicInfoByID(id uint) (*response.UserInfoForInvoice, *e.Error) {
	var user *response.UserInfoForInvoice
	query := repo.DB.Raw(`
		SELECT 
			"firstName",
			"lastName",
			email,
			phone
		FROM users
		WHERE id = ?`,
		id).Scan(&user)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return user, nil
}
