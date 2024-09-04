package productrepo

import (
	"context"
	"fmt"
	"os"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
)

var (
	errNoColourVariantByThisID_400 = e.Error{StatusCode: 400, Status: "Failed", Msg: "colourVariant doesn't exist by this ID", Err: nil}
)

func (repo *ProductsRepo) AddColourVariant(req *entities.ColourVariant, file *os.File) *e.Error {

	//initiate transaction
	tx := repo.DB.Begin()

	//defer rollback if error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//add image to cloudinary
	uploadParams := uploader.UploadParams{
		Folder:    "github.com/AbdulRahimOM/shoe-mart-ecommerce/colourVariants",
		PublicID:  uuid.New().String()[:7],
		Overwrite: true,
	}

	result, err := repo.Cld.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		return e.SetError("error while uploading file to cloudinary. err: ", err, 500)
	}

	if result.Error.Message != "" {
		return e.SetError("error while uploading file to cloudinary. result.Error: "+result.Error.Message, nil, 500)
	}

	fmt.Println("req.ImageURL: ", req.ImageURL) //url printing,.. may be required for checking purposes

	//add colourVariant
	resultGorm := tx.Create(&req)
	if resultGorm.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&resultGorm.Error)
	}

	tx.Commit()
	//need update bcoz transaction is really inefficient to the rollback because one op is in cloudinary!!!!

	return nil
}

func (repo *ProductsRepo) EditColourVariant(req *entities.ColourVariant) *e.Error {
	//check if colourVariant exists
	var temp entities.ColourVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM colour_variants
		WHERE id = ?`,
		req.ID).Scan(&temp)

	if query.Error != nil {
		return e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return &errNoColourVariantByThisID_400
	}

	//update colourVariant
	query = repo.DB.Model(&entities.ColourVariant{}).Where("id = ?", req.ID).Updates(entities.ColourVariant{
		Colour: req.Colour,
	})

	if query.Error != nil {
		return e.DBQueryError_500(&query.Error)
	}

	return nil
}

func (repo *ProductsRepo) DoColourVariantExistByAttributes(req *entities.ColourVariant) (bool, *e.Error) {

	var temp entities.ColourVariant
	query := repo.DB.Raw(`
		SELECT *
		FROM colour_variants
		WHERE "colour" = ? AND "modelId" = ?`,
		req.Colour, req.ModelID).Scan(&temp)

	if query.Error != nil {
		return false, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *ProductsRepo) GetColourVariantsUnderModel(modelID uint) (*[]entities.ColourVariant, *e.Error) {
	var colourVariants []entities.ColourVariant
	query := repo.DB.
		Preload("FkModel.FkBrand").
		Preload("FkModel.FkCategory").
		Where("\"modelId\" = ?", modelID).Find(&colourVariants)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &colourVariants, nil
}
