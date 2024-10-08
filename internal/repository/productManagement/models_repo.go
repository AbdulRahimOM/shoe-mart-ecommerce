package productrepo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"gorm.io/gorm"
)

var (
	errNoModelByThisID_400 = e.Error{StatusCode: 400, Status: "Failed", Msg: "model doesn't exist by this ID", Err: nil}
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
		return e.DBQueryError_500(&result.Error)
	}
	if result.RowsAffected == 0 {
		return &errNoModelByThisID_400
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
		return false, e.DBQueryError_500(&query.Error)
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
		return e.DBQueryError_500(&result.Error)
	}

	return nil
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
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &modelsList, nil
}

// GetSellerIdOfModel
func (repo *ModelsRepo) GetSellerIdOfModel(id uint) (uint, *e.Error) {
	var temp entities.Models
	//preloading brand and category
	query := repo.DB.Preload("FkBrand").
		Where("id = ?", id).Find(&temp)
	if query.Error != nil {
		return 0, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return 0, &errNoModelByThisID_400
	}

	return temp.FkBrand.SellerID, nil
}
