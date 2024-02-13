package usecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"os"
)

type ICategoryUC interface {
	AddCategory(req *request.AddCategoryReq) error
	GetCategories() (*[]entities.Categories, error)
	EditCategory(req *request.EditCategoryReq) error
}
type IBrandsUC interface {
	AddBrand(req *request.AddBrandReq) error
	GetBrands() (*[26]entities.BrandsByAlphabet, error)
	EditBrand(req *request.EditBrandReq) error
}

type IModelsUC interface {
	AddModel(req *request.AddModelReq) error
	EditModelName(req *request.EditModelReq) error
	GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, error)
}
type IProductsUC interface {
	AddColourVariant(sellerID uint, req *request.AddColourVariantReq, file *os.File) error
	EditColourVariant(req *request.EditColourVariantReq) error
	GetColourVariantsUnderModel(modelID uint) (*[]response.ResponseColourVarient, error)

	GetProducts() (*[]response.ResponseProduct, error)

	AddDimensionalVariant(req *request.AddDimensionalVariantReq) error

	AddStock(req *request.AddStockReq) error
	EditStock(req *request.EditStockReq) error
}
