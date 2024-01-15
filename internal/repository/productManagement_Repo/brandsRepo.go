package manageProductRepository

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BrandsRepo struct {
	DB *gorm.DB
}

func NewBrandRepository(db *gorm.DB) repoInterface.IBrandsRepo {
	return &BrandsRepo{DB: db}
}

// EditBrand
func (repo *BrandsRepo) EditBrand(req *requestModels.EditBrandReq) error {
	result := repo.DB.Model(&entities.Brands{}).Where("name = ?", req.OldName).Update("name", req.NewName)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit brand. query.Error= ", result.Error, "\n----")
		return result.Error
	}
	return nil
}

func (repo *BrandsRepo) AddBrand(req *entities.Brands) error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add brand. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *BrandsRepo) DoBrandExistsByName(name string) (bool, error) {

	var temp entities.Brands
	query := repo.DB.Raw(`
		SELECT *
		FROM brands
		WHERE name = ?`,
		name).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-brand is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *BrandsRepo) GetBrands() (*[26]entities.BrandsByAlphabet, error) {
	var brands [26]entities.BrandsByAlphabet
	for i := 0; i < 26; i++ {
		brands[i].Alphabet = string(rune(65 + i))
		query := repo.DB.Raw(`
			SELECT *
			FROM brands
			WHERE name LIKE ? OR name LIKE ?`,
			brands[i].Alphabet+"%", strings.ToLower(brands[i].Alphabet)+"%").Scan(&brands[i].Brands)

		if query.Error != nil {
			fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
			return nil, query.Error
		}
	}

	return &brands, nil
}
