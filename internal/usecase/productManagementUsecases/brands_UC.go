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

type BrandsUC struct {
	BrandsRepo repoInterface.IBrandsRepo
}

func NewBrandUseCase(repo repoInterface.IBrandsRepo) usecase.IBrandsUC {
	return &BrandsUC{BrandsRepo: repo}
}

func (uc *BrandsUC) AddBrand(req *request.AddBrandReq) *e.Error {
	//check if the brand already exists
	doBrandExistsByName, err := uc.BrandsRepo.DoBrandExistsByName(req.Name)
	if err != nil {
		return err
	}
	if doBrandExistsByName {
		return e.SetError("brand already exists", nil, 400)
	}

	var brand entities.Brands
	if err := copier.Copy(&brand, &req); err != nil {
		return &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to brand entity"), StatusCode: 500}
	}

	//add brand
	return uc.BrandsRepo.AddBrand(&brand)
}

// EditBrand
func (uc *BrandsUC) EditBrand(req *request.EditBrandReq) *e.Error {

	//check if the brand really exists
	DoBrandExistsByName, err := uc.BrandsRepo.DoBrandExistsByName(req.OldName)
	if err != nil {
		return err
	}
	if !DoBrandExistsByName {
		return e.SetError("brand doesn't exist", nil, 400)
	}

	//check if the new name already exists for another brand
	if DoBrandExistsByName, err := uc.BrandsRepo.DoBrandExistsByName(req.NewName); err != nil {
		return err
	} else if DoBrandExistsByName {
		return e.SetError("brand already exists with the sugested new name", nil, 400)
	}

	//edit brand
	return uc.BrandsRepo.EditBrand(req)
}

func (uc *BrandsUC) GetBrands() (*[26]entities.BrandsByAlphabet, *e.Error) {
	return uc.BrandsRepo.GetBrands()
}
