package productManagementHandlers

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

type ProductHandler struct {
	productUseCase usecaseInterface.IProductsUC
}

func NewProductHandler(productUseCase usecaseInterface.IProductsUC) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

// get products handler
func (h *ProductHandler) GetProducts(c *gin.Context) {
	fmt.Println("Handler ::: get products handler")

	//get products
	var products *[]response.ResponseProduct
	products, err := h.productUseCase.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting products. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SMED{
		Status:  "success",
		Message: "Products fetched successfully",
		Data:    products,
	})

}

// add stock handler
func (h *ProductHandler) AddStock(c *gin.Context) {
	fmt.Println("Enterred add stock handler")

	var req requestModels.AddStockReq
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
	//check if sellerID in token and request body match
	sellerID, err := tools.GetSellerID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding stock. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if sellerID != req.SellerID {
		fmt.Println("Seller ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "Seller ID in token and request body do not match",
		})
		return
	}

	//add stock
	if err := h.productUseCase.AddStock(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding stock. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Stock added successfully",
	})
}

// add stock handler
func (h *ProductHandler) EditStock(c *gin.Context) {
	fmt.Println("Enterred edit stock handler")

	var req requestModels.EditStockReq
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
	//check if sellerID in token and request body match
	sellerID, err := tools.GetSellerID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing stock. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if sellerID != req.SellerID {
		fmt.Println("Seller ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "Seller ID in token and request body do not match",
		})
		return
	}

	//add stock
	if err := h.productUseCase.EditStock(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing stock. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Stock edited successfully",
	})
}
