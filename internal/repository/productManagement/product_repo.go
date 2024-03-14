package productrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"

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

// GetProducts
func (repo *ProductsRepo) GetProducts() (*[]entities.Product, *e.Error) {
	var products []entities.Product
	query := repo.DB.
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Find(&products)

	if query.Error != nil {
		return nil, e.DBQueryError(&query.Error)
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
		return e.DBQueryError(&query.Error)
	}
	if query.RowsAffected == 0 {
		return e.TextError("product does not exists by this ID", 400)
	}

	//adding new count to existing stock count
	result := repo.DB.Model(&entities.Product{}).Where("id = ?", req.ProductID).Update("stock", req.AddingStockCount+earlierStock)
	if result.Error != nil {
		return e.DBQueryError(&result.Error)
	}
	return nil
}

func (repo *ProductsRepo) EditStock(req *request.EditStockReq) *e.Error {
	result := repo.DB.Model(&entities.Product{}).Where("id = ?", req.ProductID).Update("stock", req.UpdatedStockCount)
	if result.Error != nil {
		return e.DBQueryError(&result.Error)
	}
	if result.RowsAffected == 0 {
		return e.TextError("product does not exists by this ID", 400)
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
		return 0, e.DBQueryError(&query.Error)
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
		return 0, e.DBQueryError(&query.Error)
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
		return false, e.DBQueryError(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
