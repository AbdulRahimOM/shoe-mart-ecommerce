package orderrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	repo "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
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
		fmt.Println("Error occured while creating wishlist")
		return &e.Error{Err: err,StatusCode: 500}
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
		fmt.Println("Error occured while checking if wishlist exists with same name")
		return false, &e.Error{Err: err,StatusCode: 500}
	}

	return true, nil
}

func (repo *WishListRepo) AddToWishList(productID uint, wishListID uint) *e.Error {
	//add to wishlist
	wishListItem := entities.WishListItems{ProductID: productID, WishListID: wishListID}

	err := repo.DB.Create(&wishListItem).Error
	if err != nil {
		fmt.Println("Error occured while adding to wishlist")
		return &e.Error{Err: err,StatusCode: 500}
	}

	return nil
}

func (repo *WishListRepo) RemoveFromWishList(productID uint, wishListID uint) *e.Error {
	//remove from wishlist
	err := repo.DB.Where("product_id = ? AND wish_list_id = ?", productID, wishListID).Delete(&entities.WishListItems{}).Error
	if err != nil {
		fmt.Println("Error occured while removing from wishlist")
		return &e.Error{Err: err,StatusCode: 500}
	}

	return nil
}

func (repo *WishListRepo) DoesThisWishListExistForUser(userID uint, wishListID uint) (bool, *e.Error) {
	var wishList entities.WishList
	err := repo.DB.Where("user_id = ? AND id = ?", userID, wishListID).First(&wishList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			fmt.Println("Error occured while checking if wishlist exists with same name")
			return false, &e.Error{Err: err,StatusCode: 500}
		}
	}

	return true, nil
}

func (repo *WishListRepo) IsProductInWishList(productID uint, wishListID uint) (bool, *e.Error) {
	var wishListItem entities.WishListItems
	err := repo.DB.Where("product_id = ? AND wish_list_id = ?", productID, wishListID).First(&wishListItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			fmt.Println("Error occured while checking if product is in wishlist")
			return false, &e.Error{Err: err,StatusCode: 500}
		}
	}

	return true, nil
}

// GetAllWishLists
func (repo *WishListRepo) GetAllWishLists(userID uint) (*[]entities.WishList, *e.Error) {
	var wishLists []entities.WishList
	err := repo.DB.Where("user_id = ?", userID).Find(&wishLists).Error
	if err != nil {
		fmt.Println("Error occured while getting all wishlists")
		return nil, &e.Error{Err: err,StatusCode: 500}
	}

	return &wishLists, nil
}

func (repo *WishListRepo) GetWishListByID(userID uint, wishListID uint) (*string, *[]entities.Product, *e.Error) {
	var wishListName string
	var products []entities.Product

	//get wishlist name
	err := repo.DB.Model(&entities.WishList{}).Where("id = ?", wishListID).Pluck("name", &wishListName).Error
	if err != nil {
		fmt.Println("Error occured while getting wishlist name")
		return &wishListName, &products, &e.Error{Err: err,StatusCode: 500}
	}

	//get productIDs in wishlist
	var productIDs []uint
	err = repo.DB.Model(&entities.WishListItems{}).Where("wish_list_id = ?", wishListID).Pluck("product_id", &productIDs).Error
	if err != nil {
		fmt.Println("Error occured while getting productIDs in wishlist")
		return &wishListName, &products, &e.Error{Err: err,StatusCode: 500}
	}

	//get products with productIDs
	err = repo.DB.Model(&entities.Product{}).Where("id IN ?", productIDs).
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Find(&products).Error
	if err != nil {
		fmt.Println("Error occured while getting products in wishlist")
		return &wishListName, &products, &e.Error{Err: err,StatusCode: 500}
	}

	return &wishListName, &products, nil
}
