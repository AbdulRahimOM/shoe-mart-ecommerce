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

var(
	errUserIDNotExisting_404=e.Error{Status: "failed", Msg: "User ID not existing", Err: nil, StatusCode: 404}
	errSellerIDNotExisting_404=e.Error{Status: "failed", Msg: "Seller ID not existing", Err: nil, StatusCode: 404}
)

func NewAdminRepository(db *gorm.DB) repoInterface.IAdminRepo {
	return &AdminRepo{DB: db}
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

func (repo *AdminRepo) UpdateUserStatus(userID uint, newStatus string) *e.Error {
	var user entities.User
	err := repo.DB.Model(&user).Where("id = ?", userID).Update("status", newStatus).Error
	if err != nil {
		return e.DBQueryError_500(&err)
	}

	if repo.DB.RowsAffected == 0 {
		return &errUserIDNotExisting_404
	}

	return nil
}

func (repo *AdminRepo) UpdateSellerStatus(sellerID uint, newStatus string) *e.Error {
	var seller entities.Seller
	err := repo.DB.Model(&seller).Where("id = ?", sellerID).Update("status", newStatus).Error
	if err != nil {
		// return &e.Error{Status:"failed",Msg:"db query err",Err: err, StatusCode: 500}
		return e.DBQueryError_500(&err)
	}

	if repo.DB.RowsAffected == 0 {
		return &errSellerIDNotExisting_404
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
