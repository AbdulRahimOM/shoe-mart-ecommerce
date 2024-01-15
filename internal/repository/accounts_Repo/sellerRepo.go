package accRepository

import (
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type SellerRepo struct {
	DB *gorm.DB
}

func NewSellerRepository(db *gorm.DB) repoInterface.ISellerRepo {
	return &SellerRepo{DB: db}
}

func (repo *SellerRepo) GetPasswordAndSellerDetailsByEmail(email string) (string, entities.SellerDetails, error) {
	//getting password
	var hashedPassword string
	query := repo.DB.Raw(`
	SELECT password 
	FROM sellers 
	WHERE email = ?`,
		email).Scan(&hashedPassword)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", entities.SellerDetails{}, query.Error
	}

	//getting other sellerdetails
	var sellerDetails entities.SellerDetails
	query = repo.DB.Raw(`
	SELECT * 
	FROM sellers 
	WHERE email = ?`,
		email).Scan(&sellerDetails)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", entities.SellerDetails{}, query.Error
	}

	return hashedPassword, sellerDetails, nil
}

func (repo *SellerRepo) IsEmailRegistered(email string) (bool, error) {
	// fmt.Println("at repo: email=", email)
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

func (repo *SellerRepo) CreateSeller(seller *entities.Seller) error {
	sellerCreation := repo.DB.Create(&seller)
	if sellerCreation.Error != nil {
		fmt.Println("error occured while creating seller in record. ")
		return sellerCreation.Error
	}
	return nil
}
