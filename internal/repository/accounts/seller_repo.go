package accountrepo

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

func (repo *SellerRepo) GetSellerWithPwByEmail(email string) (*entities.Seller, error) {

	var seller entities.Seller
	query := repo.DB.Raw(`
	SELECT * 
	FROM sellers 
	WHERE email = ?`,
		email).Scan(&seller) //update required#1  also look  above

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &seller, nil
}

func (repo *SellerRepo) IsEmailRegistered(email string) (bool, error) {

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
