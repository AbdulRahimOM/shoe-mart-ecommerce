package manageProductRepository

import (
	"MyShoo/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
)

func (repo *ProductsRepo) AddColourVariant(req *entities.ColourVariant, file *os.File) error {

	//initiate transaction
	tx := repo.DB.Begin()

	//defer rollback if error
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("-------\npanic happened. r= ", r, "\n----")
			tx.Rollback()
		}
	}()

	//add image to cloudinary
	uploadParams := uploader.UploadParams{
		Folder:    "MyShoo/colourVariants",
		PublicID:  uuid.New().String()[:7],
		Overwrite: true,
	}

	result, err := repo.Cld.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		return errors.New("error while uploading file to cloudinary. err: " + err.Error())
	}

	if result.Error.Message != "" {
		return errors.New("error while uploading file to cloudinary. result.Error: " + result.Error.Message)
	}

	fmt.Println("req.ImageURL: ", req.ImageURL) //url printing,.. may be required for checking purposes

	//add colourVariant
	resultGorm := tx.Create(&req)
	if resultGorm.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add colourVariant. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return resultGorm.Error
	}

	tx.Commit()
	//need update bcoz transaction is really inefficient to the rollback because one op is in cloudinary!!!!

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
