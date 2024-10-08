package productrepo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repoInterface.ICategoryRepo {
	return &CategoryRepo{DB: db}
}

func (repo *CategoryRepo) AddCategory(req *entities.Categories) *e.Error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

func (repo *CategoryRepo) DoCategoryExistsByName(name string) (bool, *e.Error) {

	var temp entities.Categories
	query := repo.DB.Raw(`
        SELECT * 
        FROM categories 
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

func (repo *CategoryRepo) GetCategories() (*[]entities.Categories, *e.Error) {
	var categories []entities.Categories
	// fmt.Println("entered get categories repo")

	query := repo.DB.Raw(`
			SELECT *
			FROM categories`,
	).Scan(&categories)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &categories, nil
}

func (repo *CategoryRepo) EditCategory(req *request.EditCategoryReq) *e.Error {
	result := repo.DB.Model(&entities.Categories{}).Where("name = ?", req.OldName).Update("name", req.NewName)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}
