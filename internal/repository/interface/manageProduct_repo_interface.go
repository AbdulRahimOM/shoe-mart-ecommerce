package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	"os"
)

type ICategoryRepo interface {
	AddCategory(req *entities.Categories) error
	DoCategoryExistsByName(name string) (bool, error)
	GetCategories() (*[]entities.Categories, error)
	EditCategory(req *requestModels.EditCategoryReq) error
}
type IBrandsRepo interface {
	AddBrand(req *entities.Brands) error
	DoBrandExistsByName(name string) (bool, error)
	GetBrands() (*[26]entities.BrandsByAlphabet, error)
	EditBrand(req *requestModels.EditBrandReq) error
}

type IModelsRepo interface {
	AddModel(req *entities.Models) error
	DoModelExistsbyName(name string) (bool, error)
	EditModel(req *requestModels.EditModelReq) error
	DoModelExistsByID(id uint) (bool, error)
	GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, error)
	DoModelExistByIDAndBelongsToUser(id uint, sellerID uint) (bool, bool, error)
}

type IProductsRepo interface {
	DoColourVariantExists(req *entities.ColourVariant) (bool, error)
	AddColourVariant(req *entities.ColourVariant, file *os.File) error
	EditColourVariant(req *entities.ColourVariant) error
	DoColourVariantExistByID(id uint) (bool, error)
	GetColourVariantsUnderModel(modelID uint) (*[]entities.ColourVariant, error)

	GetProducts() (*[]entities.Product, error)

	DoDimensionalVariantExistsByAttributes(req *entities.DimensionalVariant) (bool, error)
	DoDimensionalVariantExistByID(id uint) (bool, error)

	AddDimensionalVariantAndProductCombinations(dimensionalVariant *entities.DimensionalVariant) error

	GetStockOfProduct(productID uint) (uint, error)
	AddStock(req *requestModels.AddStockReq) error
	DoProductExistsByID(id uint) (bool,error)
	EditStock(req *requestModels.EditStockReq) error
	DoesProductExistByID(id uint) (bool, error)

	GetPriceOfProduct(productID uint) (float32, error)
}
