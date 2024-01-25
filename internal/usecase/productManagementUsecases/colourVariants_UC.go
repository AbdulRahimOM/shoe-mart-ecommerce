package prodManageUsecase

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
)

func (uc *ProductsUC) AddColourVariant(sellerID uint, req *requestModels.AddColourVariantReq, file *os.File) error {
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
		fmt.Println("colourVariant already exists")
		return errors.New("colourVariant already exists")
	}

	//check if modelID exists and belongs to the seller
	doModelExists,doModelBelongsToSeller, err := uc.ModelsRepo.DoModelExistByIDAndBelongsToUser(req.ModelID, sellerID)
	if err != nil {
		fmt.Println("Error occured while checking if model exists")
		return err
	}
	if !doModelExists {
		fmt.Println("model doesn't exist with this id")
		return errors.New("model doesn't exist")
	}
	if !doModelBelongsToSeller {
		fmt.Println("model doesn't belong to this seller")
		return errors.New("model doesn't belong to this seller")
	}

	//add colourVariant
	err = uc.ProductsRepo.AddColourVariant(&colourVariant, file)
	if err != nil {
		return err
	}

	return nil
}

// EditColourVariant
func (uc *ProductsUC) EditColourVariant(req *requestModels.EditColourVariantReq) error {

	//check if the colourVariant really exists
	doColourVariantExists, err := uc.ProductsRepo.DoColourVariantExistByID(req.ID)
	if err != nil {
		fmt.Println("Error occured while checking if colourVariant exists")
		return err
	}
	if !doColourVariantExists {
		fmt.Println("colourVariant doesn't exist with this id")
		return errors.New("colourVariant doesn't exist")
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
		fmt.Println("another colourVariant already exists in these attributes")
		return errors.New("colourVariant already exists")
	}

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
