package productusecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	myMath "MyShoo/pkg/math"
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
)

func (uc *ProductsUC) AddColourVariant(sellerID uint, req *request.AddColourVariantReq, file *os.File) error {
	var colourVariant entities.ColourVariant
	if err := copier.Copy(&colourVariant, &req); err != nil {
		return err
	}

	//check if the colourVariant already exists
	doColourVariantExists, err := uc.ProductsRepo.DoColourVariantExists(&colourVariant)
	if err != nil {
		fmt.Println("Error occured while checking if colourVariant exists")
		return err
	}
	if doColourVariantExists {
		return errors.New("colourVariant already exists")
	}

	//check if modelID exists and belongs to the seller
	doModelExists, doModelBelongsToSeller, err := uc.ModelsRepo.DoModelExistByIDAndBelongsToUser(req.ModelID, sellerID)
	if err != nil {
		fmt.Println("Error occured while checking if model exists")
		return err
	}
	if !doModelExists {
		return errors.New("model doesn't exist with this id")
	}
	if !doModelBelongsToSeller {
		return errors.New("model doesn't belong to this seller")
	}

	//round off MRP and SalePrice to 2 decimal places
	colourVariant.MRP = myMath.RoundFloat32(colourVariant.MRP, 2)
	colourVariant.SalePrice = myMath.RoundFloat32(colourVariant.SalePrice, 2)

	//add colourVariant
	err = uc.ProductsRepo.AddColourVariant(&colourVariant, file)
	if err != nil {
		return err
	}

	return nil
}

// EditColourVariant
func (uc *ProductsUC) EditColourVariant(req *request.EditColourVariantReq) error {

	//check if the colourVariant really exists
	doColourVariantExists, err := uc.ProductsRepo.DoColourVariantExistByID(req.ID)
	if err != nil {
		fmt.Println("Error occured while checking if colourVariant exists")
		return err
	}
	if !doColourVariantExists {
		return errors.New("colourVariant doesn't exist with this id")
	}

	var colourVariant entities.ColourVariant
	if err := copier.Copy(&colourVariant, &req); err != nil {
		return err
	}
	//check if the coulourVariant already exists by attributes
	doColourVariantExists, err = uc.ProductsRepo.DoColourVariantExists(&colourVariant)
	if err != nil {
		fmt.Println("Error occured while checking if colourVariant exists")
		return err
	}
	if doColourVariantExists {
		return errors.New("colourVariant already exists in these attributes")
	}

	//round off MRP and SalePrice to 2 decimal places
	colourVariant.MRP = myMath.RoundFloat32(colourVariant.MRP, 2)
	colourVariant.SalePrice = myMath.RoundFloat32(colourVariant.SalePrice, 2)

	//edit colourVariant
	err = uc.ProductsRepo.EditColourVariant(&colourVariant)
	if err != nil {
		fmt.Println("Error occured while editing colourVariant")
		return err
	}

	return nil
}

func (uc *ProductsUC) GetColourVariantsUnderModel(modelID uint) (*[]response.ResponseColourVarient, error) {
	colourVariants, err := uc.ProductsRepo.GetColourVariantsUnderModel(modelID)
	if err != nil {
		fmt.Println("Error occured while getting colourVariants")
		return nil, err
	}

	//convert to response model
	var colourVariantsInResponse []response.ResponseColourVarient
	if err := copier.Copy(&colourVariantsInResponse, &colourVariants); err != nil {
		return nil, err
	}

	return &colourVariantsInResponse, nil
}
