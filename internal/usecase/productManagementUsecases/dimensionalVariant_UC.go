package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	"errors"

	"github.com/jinzhu/copier"
)

func (uc *ProductsUC) AddDimensionalVariant(req *request.AddDimensionalVariantReq) *e.Error {
	var dimensionalVariant entities.DimensionalVariant
	if err := copier.Copy(&dimensionalVariant, &req); err != nil {
		return &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to dimensionalVariant entity"), StatusCode: 500}
	}
	//check if the dimensionalVariant already exists
	doDimensionalVariantExists, err := uc.ProductsRepo.DoDimensionalVariantExistsByAttributes(&dimensionalVariant)
	if err != nil {
		return err
	}
	if doDimensionalVariantExists {
		return e.TextError("dimensionalVariant already exists", 400)
	}

	//add dimensionalVariant and its product combinations
	return uc.ProductsRepo.AddDimensionalVariantAndProductCombinations(&dimensionalVariant)
}
