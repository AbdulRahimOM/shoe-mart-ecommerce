package orderusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"

	"github.com/jinzhu/copier"
)

type WishListsUseCase struct {
	wishListsRepo repoInterface.IWishListsRepo
	productRepo   repoInterface.IProductsRepo
}

var (
	errWishListNameExists         = &e.Error{Msg: "wishlist with same name exists", StatusCode: 400}
	errWishListIdBelongsToAnother = &e.Error{Msg: "no such wishlist for this user", StatusCode: 400}
	errProductIdNotExisting       = &e.Error{Msg: "product does not exist", StatusCode: 400}
	errProductAlreadyInWishList   = &e.Error{Msg: "product is already in wishlist", Err: nil, StatusCode: 400}
	// err=&e.Error{Msg:"", StatusCode:  400}
	// err=&e.Error{Msg:"", StatusCode:  400}
	// err=&e.Error{Msg:"", StatusCode:  400}
	// err=&e.Error{Msg:"", StatusCode:  400}
	// err=&e.Error{Msg:"", StatusCode:  400}
	// err=&e.Error{Msg:"", StatusCode:  400}
	// err=&e.Error{Msg:"", StatusCode:  400}

)

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
		return errWishListNameExists
	}

	var wishList entities.WishList
	wishList.UserID = userID
	wishList.Name = req.Name

	//create wishlist
	return uc.wishListsRepo.CreateWishList(userID, &wishList)
}

// AddToWishList
func (uc *WishListsUseCase) AddToWishList(userID uint, req *request.AddToWishListReq) *e.Error {
	//check if wishlist belongs to the user
	userIDFromWishList, err := uc.wishListsRepo.GetUserIDOfWishList(req.WishListID)
	if err != nil {
		return err
	}
	if userIDFromWishList != userID {
		return errWishListIdBelongsToAnother
	}

	//check if product exists
	productExists, err := uc.productRepo.DoesProductExistByID(req.ProductID)
	if err != nil {
		return err
	}
	if !productExists {
		return errProductIdNotExisting
	}

	//check if product is already in wishlist
	productInWishList, err := uc.wishListsRepo.IsProductInWishList(req.ProductID, req.WishListID)
	if err != nil {
		return err
	}
	if productInWishList {
		return errProductAlreadyInWishList
	}

	//add product to wishlist
	return uc.wishListsRepo.AddToWishList(req.ProductID, req.WishListID)
}

// RemoveFromWishList
func (uc *WishListsUseCase) RemoveFromWishList(userID uint, req *request.RemoveFromWishListReq) *e.Error {
	//check if wishlist belongs to the user
	userIDFromWishList, err := uc.wishListsRepo.GetUserIDOfWishList(req.WishListID)
	if err != nil {
		return err
	}
	if userIDFromWishList != userID {
		return e.SetError("no such wishlist for this user", nil, 400)
	}

	//check if product is in wishlist
	productInWishList, err := uc.wishListsRepo.IsProductInWishList(req.ProductID, req.WishListID)
	if err != nil {
		return err
	}
	if !productInWishList {
		return e.SetError("product is not in wishlist", nil, 400)
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
		return nil, nil, 0, err
	}

	responseProducts = make([]response.ResponseProduct2, len(*products))
	errr := copier.Copy(&responseProducts, &products)
	if errr != nil {
		return nil, nil, 0, e.SetError("error occured while copying products to responseProducts", errr, 500)
	}

	for i, product := range *products {
		responseProducts[i].MRP = product.FkDimensionalVariation.FkColourVariant.MRP
		responseProducts[i].SalePrice = product.FkDimensionalVariation.FkColourVariant.SalePrice
	}

	return wishListName, &responseProducts, len(*products), nil
}
