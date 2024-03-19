package accountrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"

	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) repoInterface.IAdminRepo {
	return &AdminRepo{DB: db}
}

func (repo *AdminRepo) IsEmailRegisteredAsSeller(email string) (bool, *e.Error) {
	var emptyStruct struct{}
	query := repo.DB.Raw(`
	SELECT *
	FROM sellers
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

func (repo *AdminRepo) IsEmailRegisteredAsUser(email string) (bool, *e.Error) {
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

func (repo *AdminRepo) GetSellersList() (*[]entities.PwMaskedSeller, *e.Error) {
	var sellersList []entities.PwMaskedSeller
	query := repo.DB.Raw(`
	SELECT *
	FROM sellers`).
		Scan(&sellersList) //update required#2

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &sellersList, nil
}

func (repo *AdminRepo) UpdateUserStatus(email string, newStatus string) *e.Error {
	var user entities.User
	err := repo.DB.Model(&user).Where("email = ?", email).Update("status", newStatus).Error
	if err != nil {
		return e.DBQueryError_500(&err)
	}

	return nil
}

func (repo *AdminRepo) UpdateSellerStatus(email string, newStatus string) *e.Error {
	var seller entities.Seller
	err := repo.DB.Model(&seller).Where("email = ?", email).Update("status", newStatus).Error
	if err != nil {
		return &e.Error{Err: err, StatusCode: 500}
	}

	return nil
}

func (repo *AdminRepo) GetUsersList() (*[]entities.UserDetails, *e.Error) {

	var usersList []entities.UserDetails
	query := repo.DB.Raw(`
		SELECT 
			id,
			"firstName",
			"lastName",
			email,
			phone,
			status
		FROM users`).
		Scan(&usersList)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &usersList, nil
}

// returns hashed password and admin details
func (repo *AdminRepo) GetPasswordAndAdminDetailsByEmail(email string) (*string, *entities.AdminDetails, *e.Error) {

	//getting password
	var hashedPassword string
	query := repo.DB.Raw(`
	SELECT password 
	FROM admins 
	WHERE email = ?`,
		email).Scan(&hashedPassword)

	if query.Error != nil {
		return nil, nil, e.DBQueryError_500(&query.Error)
	}

	//getting other admindetails
	var adminDetails entities.AdminDetails
	query = repo.DB.Raw(`
	SELECT * 
	FROM admins 
	WHERE email = ?`,
		email).Scan(&adminDetails)

	if query.Error != nil {
		return nil, nil,  e.DBQueryError_500(&query.Error)
	}

	return &hashedPassword, &adminDetails, nil
}

func (repo *AdminRepo) IsEmailRegisteredAsAdmin(email string) (bool, *e.Error) {
	var emptyStruct struct{}
	query := repo.DB.Raw(`
        SELECT * 
        FROM admins 
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

func (repo *AdminRepo) CreateAdmin(admin *entities.Admin) *e.Error {
	result := repo.DB.Create(&admin)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	return nil
}

// VerifySeller
func (repo *AdminRepo) VerifySeller(id uint) *e.Error {
	err := repo.DB.Model(&entities.Seller{}).Where("id = ?", id).Update("status", "verified").Error
	if err != nil {
		return e.DBQueryError_500(&err)
	}

	return nil
}

// IsSellerVerified
func (repo *AdminRepo) IsSellerVerified(id uint) (bool, *e.Error) {
	var seller entities.Seller
	query := repo.DB.Raw(`
	SELECT status 
	FROM sellers 
	WHERE id = ?`,
		id).Scan(&seller)

	if query.Error != nil {
		return false, e.DBQueryError_500(&query.Error)
	}

	if seller.Status == "verified" {
		return true, nil
	}

	return false, nil
}
