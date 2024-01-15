package manageProductRepository

import (
	"MyShoo/internal/domain/entities"
	"fmt"
)

func (repo *ProductsRepo) AddColourVariant(req *entities.ColourVariant) error {
	result := repo.DB.Create(&req)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add colourVariant. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *ProductsRepo) DoColourVariantExistByID(id uint) (bool, error) {
	var temp entities.ColourVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM colour_variants
		WHERE id = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-colourVariant is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

func (repo *ProductsRepo) EditColourVariant(req *entities.ColourVariant) error {
	//check if colourVariant exists
	var temp entities.ColourVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM colour_variants
		WHERE id = ?`,
		req.ID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-colourVariant is existing or not. query.Error= ", query.Error, "\n----")
		return query.Error
	}

	if query.RowsAffected == 0 {
		return fmt.Errorf("colourVariant doesn't exist")
	}

	//update colourVariant
	query = repo.DB.Model(&entities.ColourVariant{}).Where("id = ?", req.ID).Updates(entities.ColourVariant{
		Colour: req.Colour,
	})

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't update colourVariant. query.Error= ", query.Error, "\n----")
		return query.Error
	}

	return nil
}

func (repo *ProductsRepo) DoColourVariantExists(req *entities.ColourVariant) (bool, error) {

	var temp entities.ColourVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM colour_variants
		WHERE "colour" = ? AND "modelId" = ?`,
		req.Colour, req.ModelID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-colourVariant is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		fmt.Println("rowsaffected!=0")
		return true, nil
	}
}

func (repo *ProductsRepo) GetColourVariantsUnderModel(modelID uint) (*[]entities.ColourVariant, error) {
	var colourVariants []entities.ColourVariant
	query := repo.DB.
		Preload("FkModel.FkBrand").
		Preload("FkModel.FkCategory").
		Where("\"modelId\" = ?", modelID).Find(&colourVariants)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &colourVariants, nil
}