package productrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"github.com/cloudinary/cloudinary-go"
	"gorm.io/gorm"
)

type ProductsRepo struct {
	DB  *gorm.DB
	Cld *cloudinary.Cloudinary
}

func NewProductRepository(db *gorm.DB, cloudinary *cloudinary.Cloudinary) repoInterface.IProductsRepo {
	return &ProductsRepo{
		DB:  db,
		Cld: cloudinary,
	}
}

func (repo *ProductsRepo) DoProductExistsByID(id uint) (bool, *e.Error) {
	var temp entities.Models
	query := repo.DB.Raw(`
		SELECT *
		FROM product
		WHERE "id" = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing or not. query.Error= ", query.Error, "\n----")
		return false,&e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// GetProducts
func (repo *ProductsRepo) GetProducts() (*[]entities.Product, *e.Error) {
	var products []entities.Product
	query := repo.DB.
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Find(&products)

	if query.Error != nil {
		return nil,&e.Error{Err: query.Error, StatusCode: 500}
	}

	return &products, nil
}

func (repo *ProductsRepo) AddStock(req *request.AddStockReq) *e.Error {
	//getting earlier stock
	var earlierStock uint
	query := repo.DB.Raw(`
	SELECT stock
	FROM product
	WHERE id=?`, req.ProductID).
		Scan(&earlierStock)
	if query.Error != nil {
		return&e.Error{Err: query.Error, StatusCode: 500}
	}

	//adding new count to existing stock count
	result := repo.DB.Model(&entities.Product{}).Where("id = ?", req.ProductID).Update("stock", req.AddingStockCount+earlierStock)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit brand. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}
	return nil
}

func (repo *ProductsRepo) EditStock(req *request.EditStockReq) *e.Error {
	result := repo.DB.Model(&entities.Product{}).Where("id = ?", req.ProductID).Update("stock", req.UpdatedStockCount)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't edit brand. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}
	return nil
}

func (repo *ProductsRepo) GetStockOfProduct(productID uint) (uint, *e.Error) {
	var stock uint
	query := repo.DB.Raw(`
	SELECT stock
	FROM product
	WHERE id=?`, productID).
		Scan(&stock)

	if query.Error != nil {
		return 0,&e.Error{Err: query.Error, StatusCode: 500}
	}
	return stock, nil
}

func (repo *ProductsRepo) GetPriceOfProduct(productID uint) (float32, *e.Error) {
	var product entities.Product
	//preload
	query := repo.DB.
		Preload("FkDimensionalVariation.FkColourVariant").
		Where("id = ?", productID).Find(&product)

	if query.Error != nil {
		return 0,&e.Error{Err: query.Error, StatusCode: 500}
	}
	var price float32 = product.FkDimensionalVariation.FkColourVariant.SalePrice
	return price, nil
}

func (repo *ProductsRepo) DoesProductExistByID(id uint) (bool, *e.Error) {
	var temp entities.Product
	query := repo.DB.Raw(`
		SELECT *
		FROM product
		WHERE "id" = ?`,
		id).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-model is existing or not. query.Error= ", query.Error, "\n----")
		return false,&e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
