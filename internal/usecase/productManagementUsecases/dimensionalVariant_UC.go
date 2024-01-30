package prodManageUsecase

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

func (uc *ProductsUC) AddDimensionalVariant(req *requestModels.AddDimensionalVariantReq) error {
	var dimensionalVariant entities.DimensionalVariant
	if err := copier.Copy(&dimensionalVariant, &req); err != nil {
		return err
	}
	//check if the dimensionalVariant already exists
	doDimensionalVariantExists, err := uc.ProductsRepo.DoDimensionalVariantExistsByAttributes(&dimensionalVariant)
	if err != nil {
		fmt.Println("Error occured while checking if dimensionalVariant exists")
		return err
	}
	if doDimensionalVariantExists {
		return errors.New("dimensionalVariant already exists")
	}

	//add dimensionalVariant and its product combinations
	err = uc.ProductsRepo.AddDimensionalVariantAndProductCombinations(&dimensionalVariant)
	if err != nil {
		return err
	}

	return nil
}

