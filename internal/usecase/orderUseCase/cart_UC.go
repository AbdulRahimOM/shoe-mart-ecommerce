package orderusecase

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	repo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"

	"github.com/jinzhu/copier"
)

var (
	errProductDoesNotExistInCart_404 = &e.Error{Status: "failed", Msg: "product doesn't exist in cart", Err: nil, StatusCode: 404}
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
		return e.SetError("Error while copying request to cart entity", err, 500)
	}
	//check if the product exists
	DoProductExists, quantityIfExist, err := uc.cartRepo.DoProductExistAlready(&cart)
	if err != nil {
		return err
	}
	if !DoProductExists {
		return errProductDoesNotExistInCart_404
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
		return nil, totalValue, e.SetError("Error while copying request to cart entity", err, 500)
	}

	return &responseCart, totalValue, nil
}

func (c *CartUseCase) AddToCart(req *request.AddToCartReq) *e.Error {
	var cart entities.Cart
	if err := copier.Copy(&cart, &req); err != nil {
		return e.SetError("Error while copying request to cart entity", err, 500)
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
