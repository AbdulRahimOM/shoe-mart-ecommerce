package orderrepo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	repo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"gorm.io/gorm"
)

var (
	errWishListDoesNotExist_400 = e.Error{StatusCode: 400, Status: "Failed", Msg: "WishList doesn't exist", Err: nil}
)

type WishListRepo struct {
	DB *gorm.DB
}

func NewWishListRepository(db *gorm.DB) repo.IWishListsRepo {
	return &WishListRepo{DB: db}
}

func (repo *WishListRepo) CreateWishList(userID uint, wishList *entities.WishList) *e.Error {
	err := repo.DB.Create(wishList).Error
	if err != nil {
		return e.DBQueryError_500(&err)
	}

	return nil
}

func (repo *WishListRepo) DoesWishListExistWithName(userID uint, name string) (bool, *e.Error) {
	var wishList entities.WishList
	err := repo.DB.Where("user_id = ? AND name = ?", userID, name).First(&wishList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, e.DBQueryError_500(&err)
	}

	return true, nil
}

func (repo *WishListRepo) AddToWishList(productID uint, wishListID uint) *e.Error {
	//add to wishlist
	wishListItem := entities.WishListItems{ProductID: productID, WishListID: wishListID}

	err := repo.DB.Create(&wishListItem).Error
	if err != nil {
		return e.DBQueryError_500(&err)
	}

	return nil
}

func (repo *WishListRepo) RemoveFromWishList(productID uint, wishListID uint) *e.Error {
	//remove from wishlist
	err := repo.DB.Where("product_id = ? AND wish_list_id = ?", productID, wishListID).Delete(&entities.WishListItems{}).Error
	if err != nil {
		return e.DBQueryError_500(&err)
	}

	return nil
}

func (repo *WishListRepo) GetUserIDOfWishList(wishListID uint) (uint, *e.Error) {
	//get userID of wishlist
	var userID uint
	err := repo.DB.Model(&entities.WishList{}).Where("id = ?", wishListID).Pluck("user_id", &userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &errWishListDoesNotExist_400
		} else {
			return 0, e.DBQueryError_500(&err)
		}
	}

	return userID, nil
}

func (repo *WishListRepo) IsProductInWishList(productID uint, wishListID uint) (bool, *e.Error) {
	var wishListItem entities.WishListItems
	err := repo.DB.Where("product_id = ? AND wish_list_id = ?", productID, wishListID).First(&wishListItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, e.DBQueryError_500(&err)
		}
	}

	return true, nil
}

// GetAllWishLists
func (repo *WishListRepo) GetAllWishLists(userID uint) (*[]entities.WishList, *e.Error) {
	var wishLists []entities.WishList
	err := repo.DB.Where("user_id = ?", userID).Find(&wishLists).Error
	if err != nil {
		return nil, e.DBQueryError_500(&err)
	}

	return &wishLists, nil
}

func (repo *WishListRepo) GetWishListByID(userID uint, wishListID uint) (*string, *[]entities.Product, *e.Error) {
	var wishListName string
	var products []entities.Product

	//get wishlist name
	err := repo.DB.Model(&entities.WishList{}).Where("id = ?", wishListID).Pluck("name", &wishListName).Error
	if err != nil {
		return &wishListName, &products, e.DBQueryError_500(&err)
	}

	//get productIDs in wishlist
	var productIDs []uint
	err = repo.DB.Model(&entities.WishListItems{}).Where("wish_list_id = ?", wishListID).Pluck("product_id", &productIDs).Error
	if err != nil {
		return &wishListName, &products, e.DBQueryError_500(&err)
	}

	//get products with productIDs
	err = repo.DB.Model(&entities.Product{}).Where("id IN ?", productIDs).
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Find(&products).Error
	if err != nil {
		return &wishListName, &products, e.DBQueryError_500(&err)
	}

	return &wishListName, &products, nil
}
