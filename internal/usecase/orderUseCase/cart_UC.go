package orderusecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repo "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type CartUseCase struct {
	cartRepo repo.ICartRepo
}

func NewCartUseCase(cartRepo repo.ICartRepo) usecase.ICartUC {
	return &CartUseCase{cartRepo: cartRepo}
}

func (uc *CartUseCase) DeleteFromCart(req *request.DeleteFromCartReq) error {
	var cart entities.Cart
	if err := copier.Copy(&cart, &req); err != nil {
		fmt.Println("Error occured while copying request to cart entity")
		return err
	}
	//check if the product exists
	DoProductExists, quantityIfExist, err := uc.cartRepo.DoProductExistAlready(&cart)
	if err != nil {
		fmt.Println("Error occured while checking if product exists in cart")
		return err
	}
	if !DoProductExists {
		return errors.New("product doesn't exist in cart")
	}
	if quantityIfExist == 1 {
		//delete product from cart
		err = uc.cartRepo.DeleteFromCart(req)
		if err != nil {
			fmt.Println("Error occured while deleting product from cart")
			return err
		}
	} else if quantityIfExist > 1 {
		//decrease quantity
		cart.Quantity = quantityIfExist - 1
		err = uc.cartRepo.UpdateCartItemQuantity(&cart)
		if err != nil {
			fmt.Println("Error occured while decreasing quantity")
			return err
		}
	}

	return nil
}

func (uc *CartUseCase) GetCart(userID uint) (*[]response.ResponseCartItems, float32, error) {
	var responseCart []response.ResponseCartItems

	cart, totalValue, err := uc.cartRepo.GetCart(userID)
	if err != nil {
		fmt.Println("Error occured while getting cart")
		return &responseCart, totalValue, err
	}

	if err := copier.Copy(&responseCart, &cart); err != nil {
		fmt.Println("Error occured while copying cart to response cart")
		return nil, totalValue, err
	}

	return &responseCart, totalValue, nil
}

func (c *CartUseCase) AddToCart(req *request.AddToCartReq) error {
	var cart entities.Cart
	if err := copier.Copy(&cart, &req); err != nil {
		fmt.Println("Error occured while copying request to cart entity")
		return err
	}

	//check if the product already exists
	DoProductExists, quantityIfExist, err := c.cartRepo.DoProductExistAlready(&cart)
	if err != nil {
		fmt.Println("Error occured while checking if product exists in cart")
		return err
	}
	if DoProductExists {
		cart.Quantity = quantityIfExist + 1
		//add quantity
		err = c.cartRepo.UpdateCartItemQuantity(&cart)
		if err != nil {
			fmt.Println("Error occured while adding quantity")
			return err
		}
		return nil
	} else {
		//add product to cart
		cart.Quantity = 1
		err = c.cartRepo.AddToCart(&cart)
		if err != nil {
			fmt.Println("Error occured while adding product to cart")
			return err
		}
		return nil
	}
}

// ClearCart
func (uc *CartUseCase) ClearCartOfUser(userID uint) error {
	return uc.cartRepo.ClearCartOfUser(userID)
}
