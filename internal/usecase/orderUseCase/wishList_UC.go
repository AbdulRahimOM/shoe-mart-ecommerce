package orderusecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
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
func (uc *WishListsUseCase) CreateWishList(userID uint, req *request.CreateWishListReq) error {
	//check if wishlist with same name exists
	wishListExists, err := uc.wishListsRepo.DoesWishListExistWithName(userID, req.Name)
	if err != nil {
		fmt.Println("Error occured while checking if wishlist exists with same name")
		return err
	}
	if wishListExists {
		return fmt.Errorf("wishlist with same name exists")
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
func (uc *WishListsUseCase) AddToWishList(userID uint, req *request.AddToWishListReq) error {
	//check if wishlist exists and for the user
	thisWishListExistForUser, err := uc.wishListsRepo.DoesThisWishListExistForUser(userID, req.WishListID)
	if err != nil {
		fmt.Println("Error occured while checking if wishlist exists")
		return err
	}
	if !thisWishListExistForUser {
		return fmt.Errorf("no such wishlist for this user")
	}

	//check if product exists
	productExists, err := uc.productRepo.DoesProductExistByID(req.ProductID)
	if err != nil {
		fmt.Println("Error occured while checking if product exists")
		return err
	}
	if !productExists {
		return fmt.Errorf("product does not exist")
	}

	//check if product is already in wishlist
	productInWishList, err := uc.wishListsRepo.IsProductInWishList(req.ProductID, req.WishListID)
	if err != nil {
		fmt.Println("Error occured while checking if product is already in wishlist")
		return err
	}
	if productInWishList {
		return fmt.Errorf("product is already in wishlist")
	}

	//add product to wishlist
	err = uc.wishListsRepo.AddToWishList(req.ProductID, req.WishListID)
	if err != nil {
		fmt.Println("Error occured while adding product to wishlist")
		return err
	}

	return nil
}

// RemoveFromWishList
func (uc *WishListsUseCase) RemoveFromWishList(userID uint, req *request.RemoveFromWishListReq) error {
	//check if wishlist exists and for the user
	thisWishListExistForUser, err := uc.wishListsRepo.DoesThisWishListExistForUser(userID, req.WishListID)
	if err != nil {
		fmt.Println("Error occured while checking if wishlist exists")
		return err
	}
	if !thisWishListExistForUser {
		return fmt.Errorf("no such wishlist for this user")
	}

	//check if product is in wishlist
	productInWishList, err := uc.wishListsRepo.IsProductInWishList(req.ProductID, req.WishListID)
	if err != nil {
		fmt.Println("Error occured while checking if product is already in wishlist")
		return err
	}
	if !productInWishList {
		return fmt.Errorf("product is not in wishlist")
	}

	//remove product from wishlist
	err = uc.wishListsRepo.RemoveFromWishList(req.ProductID, req.WishListID)
	if err != nil {
		fmt.Println("Error occured while removing product from wishlist")
		return err
	}

	return nil
}

// GetAllWishLists
func (uc *WishListsUseCase) GetAllWishLists(userID uint) (*[]entities.WishList, int, error) {
	//get all wishlists of user
	wishLists, err := uc.wishListsRepo.GetAllWishLists(userID)
	if err != nil {
		fmt.Println("Error occured while getting all wishlists of user")
		return nil, 0, err
	}

	return wishLists, len(*wishLists), nil
}

func (uc *WishListsUseCase) GetWishListByID(userID uint, wishListID uint) (*string, *[]response.ResponseProduct2, int, error) {
	//get wishlist
	var wishListName *string
	var responseProducts []response.ResponseProduct2
	wishListName, products, err := uc.wishListsRepo.GetWishListByID(userID, wishListID)
	if err != nil {
		fmt.Println("Error occured while getting wishlist")
		return wishListName, &responseProducts, 0, err
	}
	// Initialize responseProducts before copying
	responseProducts = make([]response.ResponseProduct2, len(*products))

	//copy products to responseProducts using copier
	err = copier.Copy(&responseProducts, &products)
	if err != nil {
		fmt.Println("Error occured while copying products to responseProducts")
		return wishListName, &responseProducts, 0, err
	}
	for i, product := range *products {
		responseProducts[i].MRP = product.FkDimensionalVariation.FkColourVariant.MRP
		responseProducts[i].SalePrice = product.FkDimensionalVariation.FkColourVariant.SalePrice
	}

	return wishListName, &responseProducts, len(*products), nil
}