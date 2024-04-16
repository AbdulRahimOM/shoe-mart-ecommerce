package orderrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	repo "MyShoo/internal/repository/interface"

	"gorm.io/gorm"
)

var(
	errNoSuchItemInCart_400=e.Error{StatusCode: 400, Status: "Failed", Msg: "No such item in cart", Err: nil}
	errCartIsEmpty_400=e.Error{StatusCode: 400, Status: "Failed", Msg: "Cart is empty", Err: nil}
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
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

func (repo *CartRepo) DeleteFromCart(req *request.DeleteFromCartReq) *e.Error {
	//check by productID and userID
	result := repo.DB.Where("\"productId\" = ? AND user_id = ?", req.ProductID, req.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	if result.RowsAffected == 0 {
		return &errNoSuchItemInCart_400
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
		return false, 0, e.DBQueryError_500(&query.Error)
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
		return nil, 0, e.DBQueryError_500(&query.Error)
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
		return e.DBQueryError_500(&result.Error)
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
		return false, e.DBQueryError_500(&query.Error)
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
		return e.DBQueryError_500(&result.Error)
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
		return 0, 0, e.DBQueryError_500(&query.Error)
	} else if queryData.TotalQuantity == 0 {
		return 0, 0, &errCartIsEmpty_400
	}
	return queryData.TotalQuantity, queryData.TotalValue, nil
}
