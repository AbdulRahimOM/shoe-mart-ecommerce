package manageProductRepository

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type ModelsRepo struct {
	DB *gorm.DB
}

func NewModelRepository(db *gorm.DB) repoInterface.IModelsRepo {
	return &ModelsRepo{DB: db}
}

// EditModelName
func (repo *ModelsRepo) EditModel(req *requestModels.EditModelReq) error {
	result := repo.DB.Model(&entities.Models{}).Where("id = ?", req.ID).Update("name", req.Name) //need update
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit model name. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

// DoModelExistsByName
func (repo *ModelsRepo) DoModelExistsbyName(name string) (bool, error) {
	var temp entities.Models
	query := repo.DB.Raw(`
		SELECT *
		FROM models
		WHERE name = ?`,
		name).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing (by name) or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// AddModel
func (repo *ModelsRepo) AddModel(req *entities.Models) error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add model. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *ModelsRepo) DoModelExistsByID(id uint) (bool, error) {//need update : is this needed?

	var temp entities.Models
	query := repo.DB.Raw(`
		SELECT *
		FROM models
		WHERE id = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing (by ID) or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *ModelsRepo) GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, error) {
	var modelsList []entities.Models
	var query *gorm.DB
	if brandExists && categoryExists {
		query = repo.DB.Preload("FkBrand").Preload("FkCategory").
			Where("\"brandId\" IN ? AND \"categoryId\" IN ?", brandIDInts, categoryIDInts).
			Find(&modelsList)
	} else if brandExists {
		query = repo.DB.Preload("FkBrand").Preload("FkCategory").
			Where("\"brandId\" IN ?", brandIDInts).
			Find(&modelsList)
	} else if categoryExists {
		query = repo.DB.Preload("FkBrand").Preload("FkCategory").
			Where("\"categoryId\" IN ?", categoryIDInts).
			Find(&modelsList)
	} else {
		query = repo.DB.Preload("FkBrand").Preload("FkCategory").
			Find(&modelsList)
	}

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get models list. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &modelsList, nil
}

// DoModelExistByIDAndBelongsToUser
func (repo *ModelsRepo) DoModelExistByIDAndBelongsToUser(id uint, sellerID uint) (bool, bool, error) {
	var temp entities.Models
	//preloading brand and category
	query := repo.DB.Preload("FkBrand").
		Where("id = ?", id).Find(&temp)
	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing by id or not. query.Error= ", query.Error, "\n----")
		return false, false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, false, nil
	} else {
		if temp.FkBrand.SellerID == sellerID {
			return true, true, nil
		} else {
			return true, false, nil
		}
	}
}
