package handlers

import (
	e "MyShoo/internal/domain/customErrors"
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	SellerUseCase usecaseInterface.ISellerUC
}

func NewSellerHandler(useCase usecaseInterface.ISellerUC) *SellerHandler {
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
// @Param req body requestModels.SellerSignUpReq{} true "Seller Sign Up Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /seller/signup [post]
func (h *SellerHandler) PostSignUp(c *gin.Context) {

	var req requestModels.SellerSignUpReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	token, err := h.SellerUseCase.SignUp(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("", err))
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Status: "success",
			Token:  *token,
		})
	}
}

// @Summary Seller Sign In Handler
// @Description Seller Sign In Handler
// @Tags Seller/Session
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body requestModels.SellerSignInReq{} true "Seller Sign In Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /seller/login [post]
func (h *SellerHandler) PostLogIn(c *gin.Context) {

	var req requestModels.SellerSignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	token, err := h.SellerUseCase.SignIn(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("", err))
		return
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Status: "success",
			Token:  *token,
		})
	}
}
