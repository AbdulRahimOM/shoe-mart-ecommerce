package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
)

type ICategoryUC interface {
	AddCategory(req *requestModels.AddCategoryReq) error
	GetCategories() (*[]entities.Categories, error)
	EditCategory(req *requestModels.EditCategoryReq) error
}
type IBrandsUC interface {
	AddBrand(req *requestModels.AddBrandReq) error
	GetBrands() (*[26]entities.BrandsByAlphabet, error)
	EditBrand(req *requestModels.EditBrandReq) error
}

type IModelsUC interface {
	AddModel(req *requestModels.AddModelReq) error
	EditModelName(req *requestModels.EditModelReq) error
	GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, error)
}
type IProductsUC interface {
	AddColourVariant(req *requestModels.AddColourVariantReq) error
	EditColourVariant(req *requestModels.EditColourVariantReq) error
	GetColourVariantsUnderModel(modelID uint) (*[]response.ResponseColourVarient, error)

	GetProducts() (*[]response.ResponseProduct, error)

	AddDimensionalVariant(req *requestModels.AddDimensionalVariantReq) error

	AddStock(req *requestModels.AddStockReq) error
	EditStock(req *requestModels.EditStockReq) error
}
