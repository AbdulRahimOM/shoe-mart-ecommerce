package productrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/tools"
	"fmt"

	"gorm.io/gorm"
)

func (repo *ProductsRepo) DoDimensionalVariantExistByID(id uint) (bool, *e.Error) {
	var temp entities.DimensionalVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM dimensional_variants
		WHERE id = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-dimensionalVariant is existing or not. query.Error= ", query.Error, "\n----")
		return false,&e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

func (repo *ProductsRepo) DoDimensionalVariantExistsByAttributes(req *entities.DimensionalVariant) (bool, *e.Error) {
	var temp entities.DimensionalVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM dimensional_variants
		WHERE "colourVariantId" = ? AND "dvIndex" = ?`,
		req.ColourVariantID, req.DVIndex).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-dimensionalVariant is existing or not. query.Error= ", query.Error, "\n----")
		return false,&e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		fmt.Println("rowsaffected!=0")
		return true, nil
	}
}

// AddDimensionalVariantAndProductCombinations(dimensionalVariant *entities.DimensionalVariant) error
func (repo *ProductsRepo) AddDimensionalVariantAndProductCombinations(dimensionalVariant *entities.DimensionalVariant) *e.Error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			fmt.Println("-------\npanic happened. couldn't add dimensionalVariant. r= ", r, "query.Error= ", result.Error, "\n----")
			tx.Rollback()
		}
	}()

	//add dimensionalVariant
	result = tx.Create(&dimensionalVariant)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add dimensionalVariant. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	//preload dimensionalVariant
	result = tx.Preload("FkColourVariant.FkModel.FkBrand").First(&dimensionalVariant, dimensionalVariant.ID)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't preload dimensionalVariant. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	//add productCombinations

	for i, size := range entities.Size {
		var productCombinations entities.Product
		productCombinations.SizeIndex = uint(i)
		// productCombinations.Stock = 0	//default value is 0, so no need //initial stock is 0
		productCombinations.DimensionalVariationID = dimensionalVariant.ID
		//skuCode,Name
		productCombinations.Name, productCombinations.SKUCode = tools.GenerateNameAndSKUCode(dimensionalVariant, &size.Size)

		//add productCombination
		result := tx.Create(&productCombinations)
		if result.Error != nil {
			fmt.Println("-------\nquery error happened. couldn't add productCombination. query.Error= ", result.Error, "\n----")
			return &e.Error{Err: result.Error, StatusCode: 500}
		}
	}

	//commit transaction
	tx.Commit()

	return nil
}
