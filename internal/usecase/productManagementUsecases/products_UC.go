package prodManageUsecase

import (
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	"fmt"
	"github.com/jinzhu/copier"
)

type ProductsUC struct {
	ProductsRepo repoInterface.IProductsRepo
	ModelsRepo  repoInterface.IModelsRepo
}

func NewProductUseCase(productRepo repoInterface.IProductsRepo, modelsRepo repoInterface.IModelsRepo) usecaseInterface.IProductsUC {
	return &ProductsUC{
		ProductsRepo: productRepo, 
		ModelsRepo: modelsRepo,
	}
}

// GetProducts
func (uc *ProductsUC) GetProducts() (*[]response.ResponseProduct, error) {
	products, err := uc.ProductsRepo.GetProducts()
	var responseProducts []response.ResponseProduct
	if err != nil {
		fmt.Println("Error occured while getting products")
		return &responseProducts, err
	}

	if err := copier.Copy(&responseProducts, &products); err != nil {
		fmt.Println("Error occured while copying products to responseProducts")
		return &responseProducts, err
	}

	return &responseProducts, nil
}

func (uc *ProductsUC) AddStock(req *requestModels.AddStockReq) error {
	//check if product exists by id
	doProductExistsByID, err := uc.ProductsRepo.DoProductExistsByID(req.ProductID)
	if err != nil {
		fmt.Println("Error while checking for product")
		return err
	}
	if !doProductExistsByID {
		fmt.Println("Product does not exists by this ID")
		return err
	}

	//add stock
	err = uc.ProductsRepo.AddStock(req)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProductsUC) EditStock(req *requestModels.EditStockReq) error {
	//check if product exists by id
	doProductExistsByID, err := uc.ProductsRepo.DoProductExistsByID(req.ProductID)
	if err != nil {
		fmt.Println("Error while checking for product")
		return err
	}
	if !doProductExistsByID {
		fmt.Println("Product does not exists by this ID")
		return err
	}

	//add stock
	err = uc.ProductsRepo.EditStock(req)
	if err != nil {
		return err
	}

	return nil
}
