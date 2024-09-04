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
	errModelAlreadyExists_409 = &e.Error{Status: "failed", Msg: "model already exists", Err: nil, StatusCode: 409}
)

type ModelsUC struct {
	ModelsRepo repoInterface.IModelsRepo
}

func NewModelUseCase(repo repoInterface.IModelsRepo) usecase.IModelsUC {
	return &ModelsUC{ModelsRepo: repo}
}

func (uc *ModelsUC) AddModel(req *request.AddModelReq) *e.Error {
	var model entities.Models
	if err := copier.Copy(&model, &req); err != nil {
		return e.SetError("Error while copying request to model entity", err, 500)
	}

	//check if the model already exists
	doModelExists, err := uc.ModelsRepo.DoModelExistsbyName(req.Name)
	if err != nil {
		return err
	}
	if doModelExists {
		return errModelAlreadyExists_409
	}

	//add model
	return uc.ModelsRepo.AddModel(&model)
}

func (uc *ModelsUC) EditModelName(req *request.EditModelReq) *e.Error {
	return uc.ModelsRepo.EditModel(req)
}

func (uc *ModelsUC) GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, *e.Error) {
	return uc.ModelsRepo.GetModelsByBrandsAndCategories(brandExists, brandIDInts, categoryExists, categoryIDInts)
}
