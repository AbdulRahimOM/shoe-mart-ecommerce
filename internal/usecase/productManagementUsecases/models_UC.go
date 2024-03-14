package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"errors"

	"github.com/jinzhu/copier"
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
		return &e.Error{Err: errors.New(err.Error() + "error occured while copying request to model entity"), StatusCode: 500}
	}

	//check if the model already exists
	doModelExists, err := uc.ModelsRepo.DoModelExistsbyName(req.Name)
	if err != nil {
		return err
	}
	if doModelExists {
		return e.TextError("model already exists", 400)
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
