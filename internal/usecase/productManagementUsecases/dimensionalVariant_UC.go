package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"

	"github.com/jinzhu/copier"
)
var (
	errDimensionalVariantAlreadyExists_409 = &e.Error{Status: "failed", Msg: "dimensionalVariant already exists", Err: nil, StatusCode: 409}
)

func (uc *ProductsUC) AddDimensionalVariant(req *request.AddDimensionalVariantReq) *e.Error {
	var dimensionalVariant entities.DimensionalVariant
	if err := copier.Copy(&dimensionalVariant, &req); err != nil {
		return e.SetError("Error while copying request to dimensionalVariant entity", err, 500)
	}
	//check if the dimensionalVariant already exists
	doDimensionalVariantExists, err := uc.ProductsRepo.DoDimensionalVariantExistsByAttributes(&dimensionalVariant)
	if err != nil {
		return err
	}
	if doDimensionalVariantExists {
		return errDimensionalVariantAlreadyExists_409
	}

	//add dimensionalVariant and its product combinations
	return uc.ProductsRepo.AddDimensionalVariantAndProductCombinations(&dimensionalVariant)
}
