package orderrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repo "MyShoo/internal/repository/interface"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) repo.ICartRepo {
	return &CartRepo{DB: db}
}

func (repo *CartRepo) AddToCart(cart *entities.Cart) *e.Error {
	result := repo.DB.Create(&cart)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add to cart. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	return nil
}

func (repo *CartRepo) DeleteFromCart(req *request.DeleteFromCartReq) *e.Error {
	//check by productID and userID
	result := repo.DB.Where("\"productId\" = ? AND user_id = ?", req.ProductID, req.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't delete from cart. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	if result.RowsAffected == 0 {
		return &e.Error{Err: errors.New("nothing deleted. no such item in cart"), StatusCode: 400}
	}

	return nil
}

func (repo *CartRepo) DoProductExistAlready(cart *entities.Cart) (bool, uint, *e.Error) {
	var temp entities.Cart
	query := repo.DB.Raw(`
		SELECT *
		FROM carts
		WHERE "productId" = ? AND user_id = ?`,
		cart.ProductID, cart.UserID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-product is existing or not. query.Error= ", query.Error, "\n----")
		return false, 0, &e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, 0, nil
	} else {
		return true, temp.Quantity, nil
	}
}

// returns cart, totalValue of cart products, error
func (repo *CartRepo) GetCart(userID uint) (*[]entities.Cart, float32, *e.Error) {
	var cart []entities.Cart
	// var totalValue float32
	query := repo.DB.
		Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
		Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
		Where("user_id = ?", userID).Find(&cart)

	if query.Error != nil {
		return nil, 0, &e.Error{Err: query.Error, StatusCode: 500}
	}
	var totalValue float32 = 0
	for i := range cart {
		totalValue += float32(cart[i].Quantity) * cart[i].FkProduct.FkDimensionalVariation.FkColourVariant.SalePrice
	}

	return &cart, totalValue, nil
}

func (repo *CartRepo) UpdateCartItemQuantity(cart *entities.Cart) *e.Error {
	result := repo.DB.Model(&entities.Cart{}).Where("\"productId\" = ? AND user_id = ?", cart.ProductID, cart.UserID).Update("quantity", cart.Quantity)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't update cart. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	return nil
}

// IsCartEmpty
func (repo *CartRepo) IsCartEmpty(userID uint) (bool, *e.Error) {
	var temp entities.Cart
	query := repo.DB.Raw(`
		SELECT *
		FROM carts
		WHERE user_id = ?`,
		userID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-cart is empty or not. query.Error= ", query.Error, "\n----")
		return false, &e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *CartRepo) ClearCartOfUser(userID uint) *e.Error {
	//delete all cart items of user where user_id = userID
	result := repo.DB.Where("user_id = ?", userID).Delete(&entities.Cart{})
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't clear cart. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	return nil
}

func (repo *CartRepo) GetQuantityAndPriceOfCart(userID uint) (uint, float32, *e.Error) {
	type data struct {
		TotalQuantity uint    `gorm:"column:totalQuantity"`
		TotalValue    float32 `gorm:"column:totalValue"`
	}
	var queryData data

	query := repo.DB.Raw(`
		SELECT 
			SUM(carts.quantity) as "totalQuantity",
			SUM(carts.quantity * colour_variants."salePrice") as "totalValue"
		FROM carts
		INNER JOIN product ON carts."productId" = product.id
		INNER JOIN dimensional_variants ON product."dimensionalVariationID" = dimensional_variants.id
		INNER JOIN colour_variants ON dimensional_variants."colourVariantId" = colour_variants.id
		WHERE user_id = ?`,
		userID).Scan(&queryData)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get quantity and price of cart. query.Error= ", query.Error, "\n----")
		return 0, 0, &e.Error{Err: query.Error, StatusCode: 500}
	} else if queryData.TotalQuantity == 0 {
		return 0, 0, &e.Error{Err: errors.New("cart is empty"), StatusCode: 400}
	}
	return queryData.TotalQuantity, queryData.TotalValue, nil
}
