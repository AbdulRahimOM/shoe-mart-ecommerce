package productrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
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
func (repo *ModelsRepo) EditModel(req *request.EditModelReq) *e.Error {
	result := repo.DB.Model(&entities.Models{}).Where("id = ?", req.ID).Update("name", req.Name) //need update
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit model name. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	return nil
}

// DoModelExistsByName
func (repo *ModelsRepo) DoModelExistsbyName(name string) (bool, *e.Error) {
	var temp entities.Models
	query := repo.DB.Raw(`
		SELECT *
		FROM models
		WHERE name = ?`,
		name).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing (by name) or not. query.Error= ", query.Error, "\n----")
		return false,&e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// AddModel
func (repo *ModelsRepo) AddModel(req *entities.Models) *e.Error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add model. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	return nil
}

func (repo *ModelsRepo) DoModelExistsByID(id uint) (bool, *e.Error) { //need update : is this needed?

	var temp entities.Models
	query := repo.DB.Raw(`
		SELECT *
		FROM models
		WHERE id = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing (by ID) or not. query.Error= ", query.Error, "\n----")
		return false,&e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *ModelsRepo) GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, *e.Error) {
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
		return nil,&e.Error{Err: query.Error, StatusCode: 500}
	}

	return &modelsList, nil
}

// DoModelExistByIDAndBelongsToUser
func (repo *ModelsRepo) DoModelExistByIDAndBelongsToUser(id uint, sellerID uint) (bool, bool, *e.Error) {
	var temp entities.Models
	//preloading brand and category
	query := repo.DB.Preload("FkBrand").
		Where("id = ?", id).Find(&temp)
	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing by id or not. query.Error= ", query.Error, "\n----")
		return false, false,&e.Error{Err: query.Error, StatusCode: 500}
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
