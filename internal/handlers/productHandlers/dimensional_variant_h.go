package producthandler

import (
	e "MyShoo/internal/domain/customErrors"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// add dimensional variant handler
// @Summary Add dimensional variant
// @Description Add dimensional variant
// @Tags Seller/Product_Management/Dimensional_Variant
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addDimensionalVariantReq body req.AddDimensionalVariantReq{} true "Add Dimensional Variant Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
func (h *ProductHandler) AddDimensionalVariant(c *gin.Context) {

	var req request.AddDimensionalVariantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//add dimensional variant
	if err := h.productUseCase.AddDimensionalVariant(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding dimensional variant. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Dimensional variant added successfully"))
}
