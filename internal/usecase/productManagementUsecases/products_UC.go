package productusecase

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"

	"github.com/jinzhu/copier"
)

type ProductsUC struct {
	ProductsRepo repoInterface.IProductsRepo
	ModelsRepo   repoInterface.IModelsRepo
}

func NewProductUseCase(productRepo repoInterface.IProductsRepo, modelsRepo repoInterface.IModelsRepo) usecase.IProductsUC {
	return &ProductsUC{
		ProductsRepo: productRepo,
		ModelsRepo:   modelsRepo,
	}
}

// GetProducts
func (uc *ProductsUC) GetProducts() (*[]response.ResponseProduct, *e.Error) {
	products, err := uc.ProductsRepo.GetProducts()
	var responseProducts []response.ResponseProduct
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&responseProducts, &products); err != nil {
		return nil, e.SetError("Error while copying products to responseProducts", err, 500)
	}

	return &responseProducts, nil
}

func (uc *ProductsUC) AddStock(req *request.AddStockReq) *e.Error {
	return uc.ProductsRepo.AddStock(req)
}

func (uc *ProductsUC) EditStock(req *request.EditStockReq) *e.Error {
	return uc.ProductsRepo.EditStock(req)
}
