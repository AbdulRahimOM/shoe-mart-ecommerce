package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"

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
		return e.SetError("error occured while copying request to category entity", err, 500)
	}

	//check if the category already exists

	DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.Name)
	if err != nil {
		return err
	}
	if DoCategoryExistsByName {
		return e.SetError("category already exists", nil, 400)
	}

	//add category
	return uc.CategoryRepo.AddCategory(&category)
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
		return e.SetError("category doesn't exist", nil, 400)
	}

	//check if the new name already exists for another category
	if DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.NewName); err != nil {
		return err
	} else if DoCategoryExistsByName {
		return e.SetError("category already exists with the sugested new name", nil, 400)
	}

	//edit category
	return uc.CategoryRepo.EditCategory(req)
}
