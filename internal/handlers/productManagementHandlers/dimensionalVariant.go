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

	var req requestModels.AddDimensionalVariantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}
	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//add dimensional variant
	if err := h.productUseCase.AddDimensionalVariant(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding dimensional variant. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSME("Dimensional variant added successfully"))
}
