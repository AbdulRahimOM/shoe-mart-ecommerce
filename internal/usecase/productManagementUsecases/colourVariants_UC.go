package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	myMath "MyShoo/pkg/math"
	"errors"
	"os"

	"github.com/jinzhu/copier"
)

var (
	errColourVariantAlreadyExistsByAttributes_409 = &e.Error{Status: "failed", Msg: "colourVariant already exists with these attributes", Err: nil, StatusCode: 409}
	errModelIDDoesNotBelongToSeller_401 = &e.Error{Status: "failed", Msg: "corrupt request", Err: errors.New("modelID does not belong to the seller"), StatusCode: 401}
	errColourVariantAlreadyExists_409 = &e.Error{Status: "failed", Msg: "colourVariant already exists with these attributes", Err: nil, StatusCode: 409}
)

func (uc *ProductsUC) AddColourVariant(sellerID uint, req *request.AddColourVariantReq, file *os.File) *e.Error {
	var colourVariant entities.ColourVariant
	if err := copier.Copy(&colourVariant, &req); err != nil {
		return e.SetError("Error while copying request to colourVariant entity", err, 500)
	}

	//check if the colourVariant already exists
	doColourVariantExists, err := uc.ProductsRepo.DoColourVariantExistByAttributes(&colourVariant)
	if err != nil {
		return err
	}
	if doColourVariantExists {
		return errColourVariantAlreadyExistsByAttributes_409
	}

	//check if modelID belongs to the seller
	sellerIDFromModel, err := uc.ModelsRepo.GetSellerIdOfModel(req.ModelID)
	if err != nil {
		return err
	}
	if sellerIDFromModel != sellerID {
		return errModelIDDoesNotBelongToSeller_401
	}

	//round off MRP and SalePrice to 2 decimal places
	colourVariant.MRP = myMath.RoundFloat32(colourVariant.MRP, 2)
	colourVariant.SalePrice = myMath.RoundFloat32(colourVariant.SalePrice, 2)

	//add colourVariant
	return uc.ProductsRepo.AddColourVariant(&colourVariant, file)
}

// EditColourVariant
func (uc *ProductsUC) EditColourVariant(req *request.EditColourVariantReq) *e.Error {

	var colourVariant entities.ColourVariant
	if err := copier.Copy(&colourVariant, &req); err != nil {
		return e.SetError("Error while copying request to colourVariant entity", err, 500)
	}
	//check if the coulourVariant already exists by attributes
	doColourVariantExists, err := uc.ProductsRepo.DoColourVariantExistByAttributes(&colourVariant)
	if err != nil {
		return err
	}
	if doColourVariantExists {
		return errColourVariantAlreadyExists_409
	}

	//round off MRP and SalePrice to 2 decimal places
	colourVariant.MRP = myMath.RoundFloat32(colourVariant.MRP, 2)
	colourVariant.SalePrice = myMath.RoundFloat32(colourVariant.SalePrice, 2)

	//edit colourVariant
	return uc.ProductsRepo.EditColourVariant(&colourVariant)
}

func (uc *ProductsUC) GetColourVariantsUnderModel(modelID uint) (*[]response.ResponseColourVarient, *e.Error) {
	colourVariants, err := uc.ProductsRepo.GetColourVariantsUnderModel(modelID)
	if err != nil {
		return nil, err
	}

	//convert to response model
	var colourVariantsInResponse []response.ResponseColourVarient
	if err := copier.Copy(&colourVariantsInResponse, &colourVariants); err != nil {
		return nil, e.SetError("Error while copying colourVariants to response model", err, 500)
	}

	return &colourVariantsInResponse, nil
}
