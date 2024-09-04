package productusecase

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"

	"github.com/jinzhu/copier"
)

var (
	errCategoryAlreadyExists_409           = &e.Error{Status: "failed", Msg: "category already exists", Err: nil, StatusCode: 409}
	errCategoryDoesNotExist_404            = &e.Error{Status: "failed", Msg: "category doesn't exist", Err: nil, StatusCode: 404}
	errCategoryAlreadyExistsInThisName_409 = &e.Error{Status: "failed", Msg: "category already exists with the sugested new name", Err: nil, StatusCode: 409}
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
		return e.SetError("Error while copying request to category entity", err, 500)
	}

	//check if the category already exists

	DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.Name)
	if err != nil {
		return err
	}
	if DoCategoryExistsByName {
		return errCategoryAlreadyExists_409
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
		return errCategoryDoesNotExist_404
	}

	//check if the new name already exists for another category
	if DoCategoryExistsByName, err := uc.CategoryRepo.DoCategoryExistsByName(req.NewName); err != nil {
		return err
	} else if DoCategoryExistsByName {
		return errCategoryAlreadyExistsInThisName_409
	}

	//edit category
	return uc.CategoryRepo.EditCategory(req)
}
