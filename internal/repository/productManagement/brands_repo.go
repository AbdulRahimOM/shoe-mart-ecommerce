package productrepo

import (
	"strings"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"gorm.io/gorm"
)

type BrandsRepo struct {
	DB *gorm.DB
}

func NewBrandRepository(db *gorm.DB) repoInterface.IBrandsRepo {
	return &BrandsRepo{DB: db}
}

// EditBrand
func (repo *BrandsRepo) EditBrand(req *request.EditBrandReq) *e.Error {
	result := repo.DB.Model(&entities.Brands{}).Where("name = ?", req.OldName).Update("name", req.NewName)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	return nil
}

func (repo *BrandsRepo) AddBrand(req *entities.Brands) *e.Error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

func (repo *BrandsRepo) DoBrandExistsByName(name string) (bool, *e.Error) {

	var temp entities.Brands
	query := repo.DB.Raw(`
		SELECT *
		FROM brands
		WHERE name = ?`,
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

func (repo *BrandsRepo) GetBrands() (*[26]entities.BrandsByAlphabet, *e.Error) {
	var brands [26]entities.BrandsByAlphabet
	for i := 0; i < 26; i++ {
		brands[i].Alphabet = string(rune(65 + i))
		query := repo.DB.Raw(`
			SELECT *
			FROM brands
			WHERE name LIKE ? OR name LIKE ?`,
			brands[i].Alphabet+"%", strings.ToLower(brands[i].Alphabet)+"%").Scan(&brands[i].Brands)

		if query.Error != nil {
			return nil, e.DBQueryError_500(&query.Error)
		}
	}

	return &brands, nil
}
