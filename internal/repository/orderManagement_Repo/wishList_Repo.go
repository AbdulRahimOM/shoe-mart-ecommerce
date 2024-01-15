package ordermanagementrepo

import (
	"MyShoo/internal/domain/entities"
	repository_interface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type WishListRepo struct {
	DB *gorm.DB
}

func NewWishListRepository(db *gorm.DB) repository_interface.IWishListsRepo {
	return &WishListRepo{DB: db}
}

func (repo *WishListRepo) CreateWishList(userID uint, wishList *entities.WishList) error {
	err := repo.DB.Create(wishList).Error
	if err != nil {
		fmt.Println("Error occured while creating wishlist")
		return err
	}

	return nil
}

func (repo *WishListRepo) DoesWishListExistWithName(userID uint, name string) (bool, error) {
	var wishList entities.WishList
	err := repo.DB.Where("user_id = ? AND name = ?", userID, name).First(&wishList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Wishlist with same name doesn't exist")
			return false, nil
		}
		fmt.Println("Error occured while checking if wishlist exists with same name")
		return false, err
	}

	fmt.Println("Wishlist with same name exists")
	return true, nil
}

func (repo *WishListRepo) AddToWishList(productID uint, wishListID uint) error {
	//add to wishlist
	wishListItem := entities.WishListItems{ProductID: productID, WishListID: wishListID}

	err := repo.DB.Create(&wishListItem).Error
	if err != nil {
		fmt.Println("Error occured while adding to wishlist")
		return err
	}

	return nil
}

func (repo *WishListRepo) RemoveFromWishList(productID uint, wishListID uint) error {
	//remove from wishlist
	err := repo.DB.Where("product_id = ? AND wish_list_id = ?", productID, wishListID).Delete(&entities.WishListItems{}).Error
	if err != nil {
		fmt.Println("Error occured while removing from wishlist")
		return err
	}

	return nil
}

func (repo *WishListRepo) DoesThisWishListExistForUser(userID uint, wishListID uint) (bool, error) {
	var wishList entities.WishList
	err := repo.DB.Where("user_id = ? AND id = ?", userID, wishListID).First(&wishList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Wishlist with same name doesn't exist")
			return false, nil
		} else {
			fmt.Println("Error occured while checking if wishlist exists with same name")
			return false, err
		}
	}

	return true, nil
}

func (repo *WishListRepo) IsProductInWishList(productID uint, wishListID uint) (bool, error) {
	var wishListItem entities.WishListItems
	err := repo.DB.Where("product_id = ? AND wish_list_id = ?", productID, wishListID).First(&wishListItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Product is not in wishlist")
			return false, nil
		} else {
			fmt.Println("Error occured while checking if product is in wishlist")
			return false, err
		}
	}

	return true, nil
}

// GetAllWishLists
func (repo *WishListRepo) GetAllWishLists(userID uint) (*[]entities.WishList, error) {
	var wishLists []entities.WishList
	err := repo.DB.Where("user_id = ?", userID).Find(&wishLists).Error
	if err != nil {
		fmt.Println("Error occured while getting all wishlists")
		return nil, err
	}

	return &wishLists, nil
}

func (repo *WishListRepo) GetWishListByID(userID uint, wishListID uint) (*string, *[]entities.Product, error) {
	var wishListName string
	var products []entities.Product

	//get wishlist name
	err := repo.DB.Model(&entities.WishList{}).Where("id = ?", wishListID).Pluck("name", &wishListName).Error
	if err != nil {
		fmt.Println("Error occured while getting wishlist name")
		return &wishListName, &products, err
	}
	
	//get productIDs in wishlist
	var productIDs []uint
	err = repo.DB.Model(&entities.WishListItems{}).Where("wish_list_id = ?", wishListID).Pluck("product_id", &productIDs).Error
	if err != nil {
		fmt.Println("Error occured while getting productIDs in wishlist")
		return &wishListName, &products, err
	}

	//get products with productIDs
	err = repo.DB.Model(&entities.Product{}).Where("id IN ?", productIDs).
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Find(&products).Error
	if err != nil {
		fmt.Println("Error occured while getting products in wishlist")
		return &wishListName, &products, err
	}

	return &wishListName, &products, nil
}
