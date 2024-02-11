package accRepository

import (
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) repoInterface.IAdminRepo {
	return &AdminRepo{DB: db}
}

func (repo *AdminRepo) IsEmailRegisteredAsSeller(email string) (bool, error) {
	var emptyStruct struct{}
	query := repo.DB.Raw(`
	SELECT *
	FROM sellers
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

func (repo *AdminRepo) IsEmailRegisteredAsUser(email string) (bool, error) {
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

func (repo *AdminRepo) GetSellersList() (*[]entities.PwMaskedSeller, error) {
	var sellersList []entities.PwMaskedSeller
	query := repo.DB.Raw(`
	SELECT *
	FROM sellers`).
		Scan(&sellersList) //update required#2

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &sellersList, nil
}

func (repo *AdminRepo) UpdateUserStatus(email string, newStatus string) error {
	var user entities.User
	err := repo.DB.Model(&user).Where("email = ?", email).Update("status", newStatus).Error
	if err != nil {
		fmt.Println("-------\nerror happened on updating status. Error= ", err, "\n----")
		return err
	}

	return nil
}

func (repo *AdminRepo) UpdateSellerStatus(email string, newStatus string) error {
	var seller entities.Seller
	err := repo.DB.Model(&seller).Where("email = ?", email).Update("status", newStatus).Error
	if err != nil {
		fmt.Println("-------\nerror happened on updating status. Error= ", err, "\n----")
		return err
	}

	return nil
}

func (repo *AdminRepo) GetUsersList() (*[]entities.UserDetails, error) {
	
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
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &usersList, nil
}

func (repo *AdminRepo) GetPasswordAndAdminDetailsByEmail(email string) (string, entities.AdminDetails, error) {

	//getting password
	var hashedPassword string
	query := repo.DB.Raw(`
	SELECT password 
	FROM admins 
	WHERE email = ?`,
		email).Scan(&hashedPassword)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", entities.AdminDetails{}, query.Error
	}

	//getting other admindetails
	var adminDetails entities.AdminDetails
	query = repo.DB.Raw(`
	SELECT * 
	FROM admins 
	WHERE email = ?`,
		email).Scan(&adminDetails)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", entities.AdminDetails{}, query.Error
	}

	return hashedPassword, adminDetails, nil
}

func (repo *AdminRepo) IsEmailRegisteredAsAdmin(email string) (bool, error) {
	fmt.Println("at repo: email=", email)
	var emptyStruct struct{}
	query := repo.DB.Raw(`
        SELECT * 
        FROM admins 
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

func (repo *AdminRepo) CreateAdmin(admin *entities.Admin) error {
	adminCreation := repo.DB.Create(&admin)
	if adminCreation.Error != nil {
		fmt.Println("error occured while creating admin in record. ")
		return adminCreation.Error
	}
	return nil
}
