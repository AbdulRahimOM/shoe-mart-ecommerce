package accounthandler

import (
	"net/http"

	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
	requestValidation "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/validation"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	SellerUseCase usecase.ISellerUC
}

func NewSellerHandler(useCase usecase.ISellerUC) *SellerHandler {
	return &SellerHandler{SellerUseCase: useCase}
}

// definefunctions those which you  have called in seller routes
// @Summary Seller Login Page
// @Description Seller Login Page
// @Tags Seller/Session
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} string
// @Router /seller/login [get]
func (h *SellerHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

// @Summary Seller Sign Up Handler
// @Description Seller Sign Up Handler
// @Tags Seller/Session
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.SellerSignUpReq{} true "Seller Sign Up Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /seller/signup [post]
func (h *SellerHandler) PostSignUp(c *gin.Context) {

	var req request.SellerSignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	token, err := h.SellerUseCase.SignUp(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Token: *token,
		})
	}
}

// @Summary Seller Sign In Handler
// @Description Seller Sign In Handler
// @Tags Seller/Session
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.SellerSignInReq{} true "Seller Sign In Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /seller/login [post]
func (h *SellerHandler) PostLogIn(c *gin.Context) {

	var req request.SellerSignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	token, err := h.SellerUseCase.SignIn(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Token: *token,
		})
	}
}
