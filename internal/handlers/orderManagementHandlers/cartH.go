package ordermanagementHandlers

import (
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUseCase usecaseInterface.ICartUC
}

func NewCartHandler(cartUseCase usecaseInterface.ICartUC) *CartHandler {
	return &CartHandler{cartUseCase: cartUseCase}
}

// add to cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	fmt.Println("Handler ::: add to cart handler")
	var req *requestModels.AddToCartReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//check if userID in token and request body match
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding to cart. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if userID != req.UserID {
		fmt.Println("User ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "User ID in token and request body do not match",
		})
		return
	}

	if err := h.cartUseCase.AddToCart(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding to cart. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Added to (or increased quantity in) cart successfully",
	})
}

// get cart
func (h *CartHandler) GetCart(c *gin.Context) {
	fmt.Println("Handler ::: get cart handler")

	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding to cart. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//get cart
	var cart *[]response.ResponseCartItems
	var totalValue float32
	cart, totalValue, err = h.cartUseCase.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting cart. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetCartResponse{
		Status:     "success",
		Message:    "Cart fetched successfully",
		Cart:       *cart,
		TotalValue: totalValue,
	})
}

// delete from cart
func (h *CartHandler) DeleteFromCart(c *gin.Context) {
	fmt.Println("Handler ::: delete from cart handler")

	var req *requestModels.DeleteFromCartReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//delete from cart
	if err := h.cartUseCase.DeleteFromCart(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error deleting from cart. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Deleted from (or decreased quantity in) cart successfully",
	})
}

// clear cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	fmt.Println("Handler ::: clear cart handler")

	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error clearing cart. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//clear cart
	if err := h.cartUseCase.ClearCartOfUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error clearing cart. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Cleared cart successfully",
	})
}
