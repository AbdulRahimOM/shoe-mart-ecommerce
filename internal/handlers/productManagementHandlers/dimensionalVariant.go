package productManagementHandlers

import (
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// add dimensional variant handler
func (h *ProductHandler) AddDimensionalVariant(c *gin.Context) {
	fmt.Println("Handler ::: add dimensional variant handler")
	var req requestModels.AddDimensionalVariantReq

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

	//add dimensional variant
	if err := h.productUseCase.AddDimensionalVariant(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding dimensional variant. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Dimensional variant added successfully",
	})
}
