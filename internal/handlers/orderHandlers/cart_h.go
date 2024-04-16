package orderhandler

import (
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecase "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUseCase usecase.ICartUC
}

func NewCartHandler(cartUseCase usecase.ICartUC) *CartHandler {
	return &CartHandler{cartUseCase: cartUseCase}
}

// add to cart
// @Summary Add to cart
// @Description Add to cart
// @Tags User/Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addToCartReq body req.AddToCartReq{} true "Add to Cart Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /cart [put]
func (h *CartHandler) AddToCart(c *gin.Context) {

	var req *request.AddToCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//check if userID in token and request body match
	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))

		return
	}
	if userID != req.UserID {
		c.JSON(http.StatusBadRequest, response.FromErrByText("user ID in token and request body do not match"))
		return
	}

	if err := h.cartUseCase.AddToCart(req); err != nil {
		c.JSON(err.StatusCode, response.MsgAndError("error adding to cart", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("added to cart successfully"))
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

	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token", errr))
		return
	}

	//get cart
	var cart *[]response.ResponseCartItems
	var totalValue float32
	cart, totalValue, err := h.cartUseCase.GetCart(userID)
	if err != nil {
		c.JSON(err.StatusCode, response.MsgAndError("error getting cart", err))
		return
	}

	c.JSON(http.StatusOK, response.GetCartResponse{
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
// @Param deleteFromCartReq body req.DeleteFromCartReq{} true "Delete from Cart Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /cart [delete]
func (h *CartHandler) DeleteFromCart(c *gin.Context) {

	var req *request.DeleteFromCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//delete from cart
	if err := h.cartUseCase.DeleteFromCart(req); err != nil {
		c.JSON(err.StatusCode, response.MsgAndError("error deleting from cart", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("deleted from cart successfully"))
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

	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//clear cart
	if err := h.cartUseCase.ClearCartOfUser(userID); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("cart cleared"))
}
