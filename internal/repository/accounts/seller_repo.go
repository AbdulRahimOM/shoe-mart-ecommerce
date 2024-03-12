package accountrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"

	"gorm.io/gorm"
)

type SellerRepo struct {
	DB *gorm.DB
}

func NewSellerRepository(db *gorm.DB) repoInterface.ISellerRepo {
	return &SellerRepo{DB: db}
}

func (repo *SellerRepo) GetSellerWithPwByEmail(email string) (*entities.Seller, *e.Error) {

	var seller entities.Seller
	query := repo.DB.Raw(`
	SELECT * 
	FROM sellers 
	WHERE email = ?`,
		email).Scan(&seller) //update required#1  also look  above

	if query.Error != nil {
		return nil, &e.Error{Err: query.Error, StatusCode: 500}
	}

	return &seller, nil
}

func (repo *SellerRepo) IsEmailRegistered(email string) (bool, *e.Error) {

	var emptyStruct struct{}
	query := repo.DB.Raw(`
        SELECT * 
        FROM sellers 
        WHERE email = ?`,
		email).Scan(&emptyStruct)

	if query.Error != nil {
		return false, &e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *SellerRepo) CreateSeller(seller *entities.Seller) *e.Error {
	sellerCreation := repo.DB.Create(&seller)
	if sellerCreation.Error != nil {
		return &e.Error{Err: sellerCreation.Error, StatusCode: 500}
	}
	return nil
}
