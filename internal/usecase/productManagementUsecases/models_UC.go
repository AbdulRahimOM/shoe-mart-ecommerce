package prodManageUsecase

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	"fmt"

	"github.com/jinzhu/copier"
)

type ModelsUC struct {
	ModelsRepo repoInterface.IModelsRepo
}

func NewModelUseCase(repo repoInterface.IModelsRepo) usecaseInterface.IModelsUC {
	return &ModelsUC{ModelsRepo: repo}
}

func (uc *ModelsUC) AddModel(req *requestModels.AddModelReq) error {
	var model entities.Models
	if err := copier.Copy(&model, &req); err != nil {
		fmt.Println("Error occured while copying request to model")
		return err
	}

	//check if the model already exists
	doModelExists, err := uc.ModelsRepo.DoModelExistsbyName(req.Name)
	if err != nil {
		fmt.Println("Error occured while checking if model exists")
		return err
	}
	if doModelExists {
		return fmt.Errorf("model already exists")
	}

	//add model
	err = uc.ModelsRepo.AddModel(&model)
	if err != nil {
		fmt.Println("Error occured while adding model")
		return err
	}
	return nil
}

func (uc *ModelsUC) EditModelName(req *requestModels.EditModelReq) error {
	//check if the model exists
	if doModelExists, err := uc.ModelsRepo.DoModelExistsByID(req.ID); err != nil {
		fmt.Println("Error occured while checking if model exists")
		return err
	}else if !doModelExists {
		return fmt.Errorf("model doesn't exist")
	}

	//edit model name
	err := uc.ModelsRepo.EditModel(req)
	if err != nil {
		fmt.Println("Error occured while editing model name")
		return err
	}

	return nil
}

func (uc *ModelsUC) GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, error) {
	var models *[]entities.Models
	models, err := uc.ModelsRepo.GetModelsByBrandsAndCategories(brandExists, brandIDInts, categoryExists, categoryIDInts)
	if err != nil {
		fmt.Println("Error occured while getting models list")
		return nil, err
	}
	return models, nil
}