package productusecase

import (
	e "MyShoo/internal/domain/customErrors"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"errors"

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
		return nil, &e.Error{Err: errors.New(err.Error() + "Error occured while copying products to responseProducts"), StatusCode: 500}
	}

	return &responseProducts, nil
}

func (uc *ProductsUC) AddStock(req *request.AddStockReq) *e.Error {
	//check if product exists by id
	doProductExistsByID, err := uc.ProductsRepo.DoProductExistsByID(req.ProductID)
	if err != nil {
		return err
	}
	if !doProductExistsByID {
		return &e.Error{Err: errors.New("product does not exists by this ID"), StatusCode: 400}
	}

	//add stock
	return uc.ProductsRepo.AddStock(req)
}

func (uc *ProductsUC) EditStock(req *request.EditStockReq) *e.Error {
	//check if product exists by id
	doProductExistsByID, err := uc.ProductsRepo.DoProductExistsByID(req.ProductID)
	if err != nil {
		return err
	}
	if !doProductExistsByID {
		return &e.Error{Err: errors.New("product does not exists by this ID"), StatusCode: 400}
	}

	//edit stock
	return uc.ProductsRepo.EditStock(req)
}
