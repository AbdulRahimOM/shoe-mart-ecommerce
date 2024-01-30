package prodManageUsecase

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type CategoryUC struct {
	CategoryRepo repoInterface.ICategoryRepo
}

func NewCategoryUseCase(categoryRepo repoInterface.ICategoryRepo) usecaseInterface.ICategoryUC {
	return &CategoryUC{
		CategoryRepo: categoryRepo,
	}
}

func (uc *CategoryUC) AddCategory(req *requestModels.AddCategoryReq) error {

	var category entities.Categories
	if err := copier.Copy(&category, &req); err != nil {
		fmt.Println("Error occured while copying request to category entity")
		return err
	}

	//check if the category already exists

	DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.Name)
	if err != nil {
		fmt.Println("Error occured while checking if category exists")
		return err
	}
	if DoCategoryExistsByName {
		return errors.New("category already exists")
	}

	//add category
	err = uc.CategoryRepo.AddCategory(&category)
	if err != nil {
		fmt.Println("Error occured while adding category")
		return err
	}

	return nil
}

func (uc *CategoryUC) GetCategories() (*[]entities.Categories, error) {
	var categories *[]entities.Categories
	categories, err := uc.CategoryRepo.GetCategories()
	if err != nil {
		fmt.Println("Error occured while getting categories list")
		return nil, err
	}
	return categories, nil
}

// edit category
func (uc *CategoryUC) EditCategory(req *requestModels.EditCategoryReq) error {
	// check if the category really exists
	DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.OldName)
	if err != nil {
		fmt.Println("Error occured while checking if category exists")
		return err
	}
	if !DoCategoryExistsByName {
		return errors.New("category doesn't exist with this old name")
	}

	//check if the new name already exists for another category
	if DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.NewName); err != nil {
		fmt.Println("Error occured while checking if category exists")
		return err
	} else if DoCategoryExistsByName {
		return errors.New("category already exists with the sugested new name")
	}

	//edit category
	err = uc.CategoryRepo.EditCategory(req)
	if err != nil {
		fmt.Println("Error occured while editing category")
		return err
	}

	return nil
}
