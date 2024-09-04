package repo

import (
	"os"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
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
	GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, *e.Error)
	GetSellerIdOfModel(id uint) (uint, *e.Error)
}

type IProductsRepo interface {
	DoColourVariantExistByAttributes(req *entities.ColourVariant) (bool, *e.Error)
	AddColourVariant(req *entities.ColourVariant, file *os.File) *e.Error
	EditColourVariant(req *entities.ColourVariant) *e.Error
	GetColourVariantsUnderModel(modelID uint) (*[]entities.ColourVariant, *e.Error)

	DoDimensionalVariantExistsByAttributes(req *entities.DimensionalVariant) (bool, *e.Error)
	AddDimensionalVariantAndProductCombinations(dimensionalVariant *entities.DimensionalVariant) *e.Error

	GetStockOfProduct(productID uint) (uint, *e.Error)
	AddStock(req *request.AddStockReq) *e.Error
	EditStock(req *request.EditStockReq) *e.Error
	DoesProductExistByID(id uint) (bool, *e.Error)

	GetProducts() (*[]entities.Product, *e.Error)
	GetPriceOfProduct(productID uint) (float32, *e.Error)
}
