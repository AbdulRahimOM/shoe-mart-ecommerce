package productrepo

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repoInterface.ICategoryRepo {
	return &CategoryRepo{DB: db}
}

func (repo *CategoryRepo) AddCategory(req *entities.Categories) error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add category. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *CategoryRepo) DoCategoryExistsByName(name string) (bool, error) {

	var temp entities.Categories
	query := repo.DB.Raw(`
        SELECT * 
        FROM categories 
        WHERE name = ?`,
		name).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-category is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *CategoryRepo) GetCategories() (*[]entities.Categories, error) {
	var categories []entities.Categories
	// fmt.Println("entered get categories repo")

	query := repo.DB.Raw(`
			SELECT *
			FROM categories`,
	).Scan(&categories)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &categories, nil
}

func (repo *CategoryRepo) EditCategory(req *request.EditCategoryReq) error {
	result := repo.DB.Model(&entities.Categories{}).Where("name = ?", req.OldName).Update("name", req.NewName)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit category. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}
