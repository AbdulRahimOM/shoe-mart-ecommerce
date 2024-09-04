package accountrepo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

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
		return nil, e.DBQueryError_500(&query.Error)
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
		return false, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *SellerRepo) CreateSeller(seller *entities.Seller) *e.Error {
	result := repo.DB.Create(&seller)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	return nil
}
