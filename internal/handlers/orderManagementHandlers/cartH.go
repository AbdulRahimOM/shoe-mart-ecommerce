package ordermanagementHandlers

import (
	e "MyShoo/internal/domain/customErrors"
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
// @Summary Add to cart
// @Description Add to cart
// @Tags User/Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addToCartReq body requestModels.AddToCartReq{} true "Add to Cart Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /cart [put]
func (h *CartHandler) AddToCart(c *gin.Context) {

	var req *requestModels.AddToCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
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

	c.JSON(http.StatusOK, response.SM{
		Status:  "success",
		Message: "Added to (or increased quantity in) cart successfully",
	})
}

// get cart
// @Summary Get cart
// @Description Get cart
// @Tags User/Cart
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetCartResponse{}
// @Failure 400 {object} response.SME{}
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {

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
// @Summary Delete from cart
// @Description Delete from cart
// @Tags User/Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param deleteFromCartReq body requestModels.DeleteFromCartReq{} true "Delete from Cart Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /cart [delete]
func (h *CartHandler) DeleteFromCart(c *gin.Context) {

	var req *requestModels.DeleteFromCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
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

	c.JSON(http.StatusOK, response.SM{
		Status:  "success",
		Message: "Deleted from (or decreased quantity in) cart successfully",
	})
}

// clear cart
// @Summary Clear cart
// @Description Clear cart
// @Tags User/Cart
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /clearcart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {

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

	c.JSON(http.StatusOK, response.SM{
		Status:  "success",
		Message: "Cleared cart successfully",
	})
}
