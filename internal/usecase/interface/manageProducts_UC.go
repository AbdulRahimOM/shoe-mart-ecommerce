package usecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"os"
)

type ICategoryUC interface {
	AddCategory(req *request.AddCategoryReq) *e.Error
	GetCategories() (*[]entities.Categories, *e.Error)
	EditCategory(req *request.EditCategoryReq) *e.Error
}
type IBrandsUC interface {
	AddBrand(req *request.AddBrandReq) *e.Error
	GetBrands() (*[26]entities.BrandsByAlphabet, *e.Error)
	EditBrand(req *request.EditBrandReq) *e.Error
}

type IModelsUC interface {
	AddModel(req *request.AddModelReq) *e.Error
	EditModelName(req *request.EditModelReq) *e.Error
	GetModelsByBrandsAndCategories(brandExists bool, brandIDInts []uint, categoryExists bool, categoryIDInts []uint) (*[]entities.Models, *e.Error)
}
type IProductsUC interface {
	AddColourVariant(sellerID uint, req *request.AddColourVariantReq, file *os.File) *e.Error
	EditColourVariant(req *request.EditColourVariantReq) *e.Error
	GetColourVariantsUnderModel(modelID uint) (*[]response.ResponseColourVarient, *e.Error)

	GetProducts() (*[]response.ResponseProduct, *e.Error)

	AddDimensionalVariant(req *request.AddDimensionalVariantReq) *e.Error

	AddStock(req *request.AddStockReq) *e.Error
	EditStock(req *request.EditStockReq) *e.Error
}
