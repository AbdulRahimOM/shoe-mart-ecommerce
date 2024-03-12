package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type CategoryUC struct {
	CategoryRepo repoInterface.ICategoryRepo
}

func NewCategoryUseCase(categoryRepo repoInterface.ICategoryRepo) usecase.ICategoryUC {
	return &CategoryUC{
		CategoryRepo: categoryRepo,
	}
}

func (uc *CategoryUC) AddCategory(req *request.AddCategoryReq) *e.Error {

	var category entities.Categories
	if err := copier.Copy(&category, &req); err != nil {
		fmt.Println("Error occured while copying request to category entity")
		return &e.Error{Err: errors.New("error occured while copying request to category entity" + err.Error()), StatusCode: 500}
	}

	//check if the category already exists

	DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.Name)
	if err != nil {
		return err
	}
	if DoCategoryExistsByName {
		return &e.Error{Err: errors.New("category already exists"), StatusCode: 400}
	}

	//add category
	err = uc.CategoryRepo.AddCategory(&category)
	if err != nil {
		fmt.Println("Error occured while adding category")
		return err
	}

	return nil
}

func (uc *CategoryUC) GetCategories() (*[]entities.Categories, *e.Error) {
	return uc.CategoryRepo.GetCategories()
}

// edit category
func (uc *CategoryUC) EditCategory(req *request.EditCategoryReq) *e.Error {
	// check if the category really exists
	DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.OldName)
	if err != nil {
		return err
	}
	if !DoCategoryExistsByName {
		return &e.Error{Err: errors.New("category doesn't exist"), StatusCode: 400}
	}

	//check if the new name already exists for another category
	if DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.NewName); err != nil {
		return err
	} else if DoCategoryExistsByName {
		return &e.Error{Err: errors.New("category already exists with the sugested new name"),StatusCode: 400}
	}

	//edit category
	return uc.CategoryRepo.EditCategory(req)
}
