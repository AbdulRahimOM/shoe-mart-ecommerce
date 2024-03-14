package orderusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repo "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"errors"

	"github.com/jinzhu/copier"
)

type CartUseCase struct {
	cartRepo repo.ICartRepo
}

func NewCartUseCase(cartRepo repo.ICartRepo) usecase.ICartUC {
	return &CartUseCase{cartRepo: cartRepo}
}

func (uc *CartUseCase) DeleteFromCart(req *request.DeleteFromCartReq) *e.Error {
	var cart entities.Cart
	if err := copier.Copy(&cart, &req); err != nil {
		return &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to cart entity"), StatusCode: 500}
	}
	//check if the product exists
	DoProductExists, quantityIfExist, err := uc.cartRepo.DoProductExistAlready(&cart)
	if err != nil {
		return err
	}
	if !DoProductExists {
		return e.TextError("product doesn't exist in cart", 400)
	}
	if quantityIfExist == 1 {
		//delete product from cart
		return uc.cartRepo.DeleteFromCart(req)
	} else { //case: quantityIfExist > 1 {
		//decrease quantity
		cart.Quantity = quantityIfExist - 1
		return uc.cartRepo.UpdateCartItemQuantity(&cart)
	}
}

func (uc *CartUseCase) GetCart(userID uint) (*[]response.ResponseCartItems, float32, *e.Error) {
	var responseCart []response.ResponseCartItems

	cart, totalValue, err := uc.cartRepo.GetCart(userID)
	if err != nil {
		return &responseCart, totalValue, err
	}

	if err := copier.Copy(&responseCart, &cart); err != nil {
		return nil, totalValue, &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to cart entity"), StatusCode: 500}
	}

	return &responseCart, totalValue, nil
}

func (c *CartUseCase) AddToCart(req *request.AddToCartReq) *e.Error {
	var cart entities.Cart
	if err := copier.Copy(&cart, &req); err != nil {
		return &e.Error{Err: errors.New(err.Error() + "Error occured while copying request to cart entity"), StatusCode: 500}
	}

	//check if the product already exists
	DoProductExists, quantityIfExist, err := c.cartRepo.DoProductExistAlready(&cart)
	if err != nil {
		return err
	}
	if DoProductExists {
		cart.Quantity = quantityIfExist + 1
		//add quantity
		return c.cartRepo.UpdateCartItemQuantity(&cart)
	} else {
		//add product to cart
		cart.Quantity = 1
		return c.cartRepo.AddToCart(&cart)
	}
}

// ClearCart
func (uc *CartUseCase) ClearCartOfUser(userID uint) *e.Error {
	return uc.cartRepo.ClearCartOfUser(userID)
}
