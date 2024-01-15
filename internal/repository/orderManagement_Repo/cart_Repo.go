package ordermanagementrepo

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	repository_interface "MyShoo/internal/repository/interface"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) repository_interface.ICartRepo {
	return &CartRepo{DB: db}
}

func (repo *CartRepo) AddToCart(cart *entities.Cart) error {
	result := repo.DB.Create(&cart)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add to cart. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *CartRepo) DeleteFromCart(req *requestModels.DeleteFromCartReq) error {
	//check by productID and userID
	result := repo.DB.Where("\"productId\" = ? AND user_id = ?", req.ProductID, req.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't delete from cart. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("couldn't delete from cart")
	}

	return nil
}

func (repo *CartRepo) DoProductExistAlready(cart *entities.Cart) (bool, uint, error) {
	var temp entities.Cart
	query := repo.DB.Raw(`
		SELECT *
		FROM carts
		WHERE "productId" = ? AND user_id = ?`,
		cart.ProductID, cart.UserID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-product is existing or not. query.Error= ", query.Error, "\n----")
		return false, 0, query.Error
	}

	if query.RowsAffected == 0 {
		return false, 0, nil
	} else {
		return true, temp.Quantity, nil
	}
}

func (repo *CartRepo) GetCart(userID uint) (*[]entities.Cart, error) {
	var cart []entities.Cart
	// var totalValue float32
	query := repo.DB.
		Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Where("user_id = ?", userID).Find(&cart)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &cart, nil
}

func (repo *CartRepo) UpdateCartItemQuantity(cart *entities.Cart) error {
	result := repo.DB.Model(&entities.Cart{}).Where("\"productId\" = ? AND user_id = ?", cart.ProductID, cart.UserID).Update("quantity", cart.Quantity)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't update cart. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

// IsCartEmpty
func (repo *CartRepo) IsCartEmpty(userID uint) (bool, error) {
	var temp entities.Cart
	query := repo.DB.Raw(`
		SELECT *
		FROM carts
		WHERE user_id = ?`,
		userID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-cart is empty or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *CartRepo) ClearCartOfUser(userID uint) error {
	//delete all cart items of user where user_id = userID
	result := repo.DB.Where("user_id = ?", userID).Delete(&entities.Cart{})
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't clear cart. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}