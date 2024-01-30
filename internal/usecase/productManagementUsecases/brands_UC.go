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

type BrandsUC struct {
	BrandsRepo repoInterface.IBrandsRepo
}

func NewBrandUseCase(repo repoInterface.IBrandsRepo) usecaseInterface.IBrandsUC {
	return &BrandsUC{BrandsRepo: repo}
}

func (uc *BrandsUC) AddBrand(req *requestModels.AddBrandReq) error {
	//check if the brand already exists
	doBrandExistsByName, err := uc.BrandsRepo.DoBrandExistsByName(req.Name)
	if err != nil {
		fmt.Println("Error occured while checking if brand exists")
		return err
	}
	if doBrandExistsByName {
		return errors.New("brand already exists")
	}
	
	var brand entities.Brands
	if err := copier.Copy(&brand, &req); err != nil {
		return err
	}

	//add brand
	err = uc.BrandsRepo.AddBrand(&brand)
	if err != nil {
		fmt.Println("Error occured while adding brand")
		return err
	}

	return nil
}

// EditBrand
func (uc *BrandsUC) EditBrand(req *requestModels.EditBrandReq) error {

	//check if the brand really exists
	DoBrandExistsByName, err := uc.BrandsRepo.DoBrandExistsByName(req.OldName)
	if err != nil {
		fmt.Println("Error occured while checking if brand exists")
		return err
	}
	if !DoBrandExistsByName {
		return errors.New("brand doesn't exist")
	}

	//check if the new name already exists for another brand
	if DoBrandExistsByName, err := uc.BrandsRepo.DoBrandExistsByName(req.NewName); err != nil {
		fmt.Println("Error occured while checking if brand exists")
		return err
	} else if DoBrandExistsByName {
		return errors.New("brand already exists with the sugested new name")
	}

	//edit brand
	err = uc.BrandsRepo.EditBrand(req)
	if err != nil {
		fmt.Println("Error occured while editing brand")
		return err
	}

	return nil
}

func (uc *BrandsUC) GetBrands() (*[26]entities.BrandsByAlphabet, error) {
	brands, err := uc.BrandsRepo.GetBrands()
	if err != nil {
		fmt.Println("Error occured while getting brands list")
		return nil, err
	}
	return brands, nil
}
