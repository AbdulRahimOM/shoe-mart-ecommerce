package accRepository

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (repo *UserRepo) GetAddressNameByID(id uint) (string, error) {
	var name string
	query := repo.DB.Raw(`
		SELECT "addressName"
		FROM user_addresses
		WHERE id = ?`,
		id).Scan(&name)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get address name by id. query.Error= ", query.Error, "\n----")
		return "", query.Error
	}

	return name, nil
}

func (repo *UserRepo) DoAddressExistsByID(id uint) (bool, error) {
	var temp entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE id = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-address is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *UserRepo) EditUserAddress(newAddress *entities.UserAddress) error {
	result := repo.DB.Model(&entities.UserAddress{}).Where("id = ?", newAddress.ID).Updates(newAddress)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit address. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *UserRepo) AddUserAddress(newAddress *entities.UserAddress) error {
	result := repo.DB.Create(&newAddress)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add to cart. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *UserRepo) DoAddressNameExists(name string) (bool, error) {
	var temp entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE "addressName" = ?`,
		name).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-address-name is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func NewUserRepository(db *gorm.DB) repoInterface.IUserRepo {
	return &UserRepo{DB: db}
}
func (repo *UserRepo) UpdateUserStatus(email string, newStatus string) error {
	var user entities.User
	err := repo.DB.Model(&user).Where("email = ?", email).Update("status", newStatus).Error
	if err != nil {
		fmt.Println("-------\nerror happened on updating status. Error= ", err, "\n----")
		return err
	}

	return nil
}

func (repo *UserRepo) GetPasswordAndUserDetailsByEmail(email string) (string, entities.UserDetails, error) {
	//getting password
	var hashedPassword string
	query := repo.DB.Raw(`
	SELECT password 
	FROM users 
	WHERE email = ?`,
		email).Scan(&hashedPassword)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", entities.UserDetails{}, query.Error
	}

	//getting other userdetails
	var userDetails entities.UserDetails
	query = repo.DB.Raw(`
	SELECT * 
	FROM users 
	WHERE email = ?`,
		email).Scan(&userDetails)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", entities.UserDetails{}, query.Error
	}

	return hashedPassword, userDetails, nil
}

func (repo *UserRepo) IsEmailRegistered(email string) (bool, error) {
	// fmt.Println("at repo: email=", email)
	var emptyStruct struct{}
	query := repo.DB.Raw(`
        SELECT * 
        FROM users 
        WHERE email = ?`,
		email).Scan(&emptyStruct)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *UserRepo) CreateUser(user *entities.User) error {
	userCreation := repo.DB.Create(&user)
	if userCreation.Error != nil {
		fmt.Println("error occured while creating user in record. ")
		return userCreation.Error
	}
	return nil
}

// DeleteUserAddress
func (repo *UserRepo) DeleteUserAddress(id uint) error {
	result := repo.DB.Delete(&entities.UserAddress{}, id)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't delete address. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

// GetUserIDFromAddressID
func (repo *UserRepo) GetUserIDFromAddressID(id uint) (uint, error) {
	var userID uint
	query := repo.DB.Raw(`
		SELECT "userId"
		FROM user_addresses
		WHERE id = ?`,
		id).Scan(&userID)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get userID from addressID. query.Error= ", query.Error, "\n----")
		return 0, query.Error
	}

	return userID, nil
}

func (repo *UserRepo) DoUserExistsByID(userId uint) (bool, error) {
	var temp entities.User
	query := repo.DB.Raw(`
		SELECT *
		FROM users
		WHERE id = ?`,
		userId).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-user is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *UserRepo) GetUserAddresses(userId uint) (*[]entities.UserAddress, error) {
	var addresses []entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE "userId" = ?`,
		userId).Scan(&addresses)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get user addresses. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &addresses, nil
}

func (repo *UserRepo) GetProfile(userID uint) (*entities.UserDetails, error) {
	var user *entities.UserDetails
	query := repo.DB.Raw(`
		SELECT *
		FROM users
		WHERE id = ?`,
		userID).Scan(&user)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get user profile. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return user, nil
}

// EditProfile implements repository_interface.IUserRepo.
func (repo *UserRepo) EditProfile(userID uint, req *requestModels.EditProfileReq) error {

	fmt.Println("from user repo: userID=", userID, "req=", req)
	result := repo.DB.Model(&entities.User{}).Where("id = ?", userID).Updates(req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit profile. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

// GetEmailByID implements repository_interface.IUserRepo.
func (repo *UserRepo) GetEmailByID(userID uint) (string, error) {
	var email string
	query := repo.DB.Raw(`
		SELECT email
		FROM users
		WHERE id = ?`,
		userID).Scan(&email)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get email by id. query.Error= ", query.Error, "\n----")
		return "", query.Error
	}

	return email, nil
}

func (repo *UserRepo) GetUserByEmail(email string)  (*entities.User, error) {
	var user entities.User
	query := repo.DB.Raw(`
		SELECT *
		FROM users
		WHERE email = ?`,
		email).Scan(&user)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get user phone by email. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &user, nil
}

func (repo *UserRepo) ResetPassword(id uint, newPassword *string) error {
	result := repo.DB.Model(&entities.User{}).Where("id = ?", id).Update("password", newPassword)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't reset password. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

//DoAddressExistsByIDForUser
func (repo *UserRepo) DoAddressExistsByIDForUser(id uint, userID uint) (bool, error) {
	var temp entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE id = ? AND "userId" = ?`,
		id, userID).Scan(&temp)
	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-address is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

//	GetUserAddress(userID uint, addressID uint) (*entities.UserAddress, error)
func (repo *UserRepo) GetUserAddress(addressID uint) (*entities.UserAddress, error) {
	var address entities.UserAddress
	query := repo.DB.Raw(`
		SELECT *
		FROM user_addresses
		WHERE id = ?`,
		addressID).Scan(&address)
	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get address. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &address, nil
}
