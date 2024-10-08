package productrepo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"github.com/cloudinary/cloudinary-go"
	"gorm.io/gorm"
)

var (
	errNoProductByThisID_400 = e.Error{StatusCode: 400, Status: "Failed", Msg: "product doesn't exist by this ID", Err: nil}
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
		return nil, e.DBQueryError_500(&query.Error)
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
		return e.DBQueryError_500(&query.Error)
	}
	if query.RowsAffected == 0 {
		return &errNoProductByThisID_400
	}

	//adding new count to existing stock count
	result := repo.DB.Model(&entities.Product{}).Where("id = ?", req.ProductID).Update("stock", req.AddingStockCount+earlierStock)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	return nil
}

func (repo *ProductsRepo) EditStock(req *request.EditStockReq) *e.Error {
	result := repo.DB.Model(&entities.Product{}).Where("id = ?", req.ProductID).Update("stock", req.UpdatedStockCount)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	if result.RowsAffected == 0 {
		return &errNoProductByThisID_400
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
		return 0, e.DBQueryError_500(&query.Error)
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
		return 0, e.DBQueryError_500(&query.Error)
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
		return false, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
