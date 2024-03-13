package orderusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type WishListsUseCase struct {
	wishListsRepo repoInterface.IWishListsRepo
	productRepo   repoInterface.IProductsRepo
}

func NewWishListUseCase(wishListsRepo repoInterface.IWishListsRepo, productRepo repoInterface.IProductsRepo) usecase.IWishListsUC {
	return &WishListsUseCase{
		wishListsRepo: wishListsRepo,
		productRepo:   productRepo,
	}
}

// CreateWishList
func (uc *WishListsUseCase) CreateWishList(userID uint, req *request.CreateWishListReq) *e.Error {
	//check if wishlist with same name exists
	wishListExists, err := uc.wishListsRepo.DoesWishListExistWithName(userID, req.Name)
	if err != nil {
		return err
	}
	if wishListExists {
		return &e.Error{Err: errors.New("wishlist with same name exists"), StatusCode: 400}
	}

	var wishList entities.WishList
	wishList.UserID = userID
	wishList.Name = req.Name

	//create wishlist
	err = uc.wishListsRepo.CreateWishList(userID, &wishList)
	if err != nil {
		fmt.Println("Error occured while creating wishlist")
		return err
	}

	return nil
}

// AddToWishList
func (uc *WishListsUseCase) AddToWishList(userID uint, req *request.AddToWishListReq) *e.Error {
	//check if wishlist exists and for the user
	thisWishListExistForUser, err := uc.wishListsRepo.DoesThisWishListExistForUser(userID, req.WishListID)
	if err != nil {
		return err
	}
	if !thisWishListExistForUser {
		return &e.Error{Err: errors.New("no such wishlist for this user"), StatusCode: 400}
	}

	//check if product exists
	productExists, err := uc.productRepo.DoesProductExistByID(req.ProductID)
	if err != nil {
		return err
	}
	if !productExists {
		return &e.Error{Err: errors.New("product does not exist"), StatusCode: 400}
	}

	//check if product is already in wishlist
	productInWishList, err := uc.wishListsRepo.IsProductInWishList(req.ProductID, req.WishListID)
	if err != nil {
		return err
	}
	if productInWishList {
		return &e.Error{Err: errors.New("product is already in wishlist"), StatusCode: 400}
	}

	//add product to wishlist
	return uc.wishListsRepo.AddToWishList(req.ProductID, req.WishListID)
}

// RemoveFromWishList
func (uc *WishListsUseCase) RemoveFromWishList(userID uint, req *request.RemoveFromWishListReq) *e.Error {
	//check if wishlist exists and for the user
	thisWishListExistForUser, err := uc.wishListsRepo.DoesThisWishListExistForUser(userID, req.WishListID)
	if err != nil {
		return err
	}
	if !thisWishListExistForUser {
		return &e.Error{Err: errors.New("no such wishlist for this user"), StatusCode: 400}
	}

	//check if product is in wishlist
	productInWishList, err := uc.wishListsRepo.IsProductInWishList(req.ProductID, req.WishListID)
	if err != nil {
		return err
	}
	if !productInWishList {
		return &e.Error{Err: errors.New("product is not in wishlist"), StatusCode: 400}
	}

	//remove product from wishlist
	return uc.wishListsRepo.RemoveFromWishList(req.ProductID, req.WishListID)
}

// GetAllWishLists
func (uc *WishListsUseCase) GetAllWishLists(userID uint) (*[]entities.WishList, int, *e.Error) {
	wishLists, err := uc.wishListsRepo.GetAllWishLists(userID)
	if err != nil {
		return nil, 0, err
	}

	return wishLists, len(*wishLists), nil
}

func (uc *WishListsUseCase) GetWishListByID(userID uint, wishListID uint) (*string, *[]response.ResponseProduct2, int, *e.Error) {
	var wishListName *string
	var responseProducts []response.ResponseProduct2
	wishListName, products, err := uc.wishListsRepo.GetWishListByID(userID, wishListID)
	if err != nil {
		return nil,nil, 0, err
	}
	
	responseProducts = make([]response.ResponseProduct2, len(*products))
	errr := copier.Copy(&responseProducts, &products)
	if errr != nil {
		return nil, nil, 0,&e.Error{Err: errors.New("error occured while copying products to responseProducts"+errr.Error()), StatusCode: 500}
	}

	for i, product := range *products {
		responseProducts[i].MRP = product.FkDimensionalVariation.FkColourVariant.MRP
		responseProducts[i].SalePrice = product.FkDimensionalVariation.FkColourVariant.SalePrice
	}

	return wishListName, &responseProducts, len(*products), nil
}
