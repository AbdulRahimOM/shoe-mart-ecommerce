package repo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	"os"
)

type ICategoryRepo interface {
	AddCategory(req *entities.Categories) *e.Error
	DoCategoryExistsByName(name string) (bool, *e.Error)
	GetCategories() (*[]entities.Categories, *e.Error)
	EditCategory(req *request.EditCategoryReq) *e.Error
}
type IBrandsRepo interface {
	AddBrand(req *entities.Brands) *e.Error
	DoBrandExistsByName(name string) (bool, *e.Error)
	GetBrands() (*[26]entities.BrandsByAlphabet, *e.Error)
	EditBrand(req *request.EditBrandReq) *e.Error
}

type IModelsRepo interface {
	AddModel(req *entities.Models) *e.Error
	DoModelExistsbyName(name string) (bool, *e.Error)
	EditModel(req *request.EditModelReq) *e.Error
	DoModelExistsByID(id uint) (bool, *e.Error)
	GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, *e.Error)
	DoModelExistByIDAndBelongsToUser(id uint, sellerID uint) (bool, bool, *e.Error)
}

type IProductsRepo interface {
	DoColourVariantExists(req *entities.ColourVariant) (bool, *e.Error)
	AddColourVariant(req *entities.ColourVariant, file *os.File) *e.Error
	EditColourVariant(req *entities.ColourVariant) *e.Error
	DoColourVariantExistByID(id uint) (bool, *e.Error)
	GetColourVariantsUnderModel(modelID uint) (*[]entities.ColourVariant, *e.Error)

	GetProducts() (*[]entities.Product, *e.Error)

	DoDimensionalVariantExistsByAttributes(req *entities.DimensionalVariant) (bool, *e.Error)
	DoDimensionalVariantExistByID(id uint) (bool, *e.Error)

	AddDimensionalVariantAndProductCombinations(dimensionalVariant *entities.DimensionalVariant) *e.Error

	GetStockOfProduct(productID uint) (uint, *e.Error)
	AddStock(req *request.AddStockReq) *e.Error
	DoProductExistsByID(id uint) (bool, *e.Error)
	EditStock(req *request.EditStockReq) *e.Error
	DoesProductExistByID(id uint) (bool, *e.Error)

	GetPriceOfProduct(productID uint) (float32, *e.Error)
}
