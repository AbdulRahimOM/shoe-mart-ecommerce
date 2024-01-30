package handlers

import (
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
// @Tags seller
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /seller/login [get]
func (h *SellerHandler) GetLogin(c *gin.Context) {
	fmt.Println("Handler ::: GET login handler")

	c.JSON(http.StatusOK, "token")
}

// @Summary Seller Sign Up Handler
// @Description Seller Sign Up Handler
// @Tags seller
// @Accept json
// @Produce json
// @Param signUpReq body requestModels.SellerSignUpReq true "Seller Sign Up Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /seller/signup [post]
func (h *SellerHandler) PostSignUp(c *gin.Context) {
	fmt.Println("=============\nentered POST sign-up handler")

	var signUpReq requestModels.SellerSignUpReq

	if err := c.Bind(&signUpReq); err != nil {
		fmt.Println("\n", "Error occured while signing up. Error while binding request"+err.Error(), "\n.")

		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error while binding request",
			Error:   err.Error(),
		})
		return
	}
	//validation
	if err := requestValidation.ValidateRequest(signUpReq); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})
		return
	}

	token, err := h.SellerUseCase.SignUp(&signUpReq)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "Error occured while signing up. Try again. Error:" + err.Error() ////////////////////////////////
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "",
			Token:   *token,
		})
	}
}

// @Summary Seller Sign In Handler
// @Description Seller Sign In Handler
// @Tags seller
// @Accept json
// @Produce json
// @Param loginReq body requestModels.SellerSignInReq true "Seller Sign In Request"
// @Success 200 {object} string
// @Failure 400 {object} string
func (h *SellerHandler) PostLogIn(c *gin.Context) {
	fmt.Println("=============\nentered \"POST login\" handler")
	// Print the raw request
	fmt.Println("c.request=", c.Request, "\n.")
	fmt.Println("c.request.body=", c.Request.Body)
	var signInReq requestModels.SellerSignInReq

	if err := c.ShouldBindJSON(&signInReq); err != nil {
		fmt.Println("\nerror binding the requewst\n.")
		errResponse := "error binding the requewst. Try again. Error:" + err.Error()
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(signInReq); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})
		return
	}

	token, err := h.SellerUseCase.SignIn(&signInReq)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "error while signing in"
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})
		return
	} else {
		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "",
			Token:   *token,
		})
	}
}
