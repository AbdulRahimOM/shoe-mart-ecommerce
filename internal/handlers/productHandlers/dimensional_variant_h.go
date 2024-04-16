package producthandler

import (
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
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
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//add dimensional variant
	if err := h.productUseCase.AddDimensionalVariant(&req); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("dimensional variant added"))
}
