package productrepo

import (
	e "MyShoo/internal/domain/customErrors"
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

func (repo *CategoryRepo) AddCategory(req *entities.Categories) *e.Error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add category. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error,StatusCode: 500}
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
		fmt.Println("-------\nquery error happened. couldn't check if-category is existing or not. query.Error= ", query.Error, "\n----")
		return false, &e.Error{Err: query.Error,StatusCode: 500}
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
		return nil, &e.Error{Err: query.Error,StatusCode: 500}
	}

	return &categories, nil
}

func (repo *CategoryRepo) EditCategory(req *request.EditCategoryReq) *e.Error {
	result := repo.DB.Model(&entities.Categories{}).Where("name = ?", req.OldName).Update("name", req.NewName)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit category. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error,StatusCode: 500}
	}

	return nil
}
