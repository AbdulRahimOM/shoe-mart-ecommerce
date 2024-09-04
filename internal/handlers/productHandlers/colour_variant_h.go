package producthandler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/tools"
	requestValidation "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/validation"

	"github.com/gin-gonic/gin"
)

// add colour variant handler
// @Summary Add colour variant
// @Description Add colour variant
// @Tags Seller/Product_Management/Colour_Variant
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addColourVariantReq body req.AddColourVariantReq{} true "Add Colour Variant Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /seller/addcolourvariant [post]
func (cvh *ProductHandler) AddColourVariant(c *gin.Context) {

	var req request.AddColourVariantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//get sellerID from token
	sellerID, errr := tools.GetSellerID(c)
	if errr != nil {
		c.JSON(http.StatusForbidden, response.MsgAndError("error getting sellerID from token:", errr))
		return
	}

	//image upload handling
	formFile, errr := c.FormFile("imageUrl")
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error getting image file from request:", errr))
		return
	}

	path := filepath.Join(os.TempDir(), formFile.Filename)
	if errr := c.SaveUploadedFile(formFile, path); errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error saving image file to temp dir:", errr))
		return
	}

	file, errr := os.Open(path)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error opening image file:", errr))
		return
	}
	defer file.Close()

	//add colour variant
	if err := cvh.productUseCase.AddColourVariant(sellerID, &req, file); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("colour variant added"))
}

// edit colour variant handler
// @Summary Edit colour variant
// @Description Edit colour variant
// @Tags Admin/Product_Management/Colour_Variant
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param editColourVariantReq body req.EditColourVariantReq{} true "Edit Colour Variant Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/editcolourvariant [patch]
func (cvh *ProductHandler) EditColourVariant(c *gin.Context) {

	var req request.EditColourVariantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//add colour variant
	if err := cvh.productUseCase.EditColourVariant(&req); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("colour variant edited"))
}

// get colour variants under model handler
// @Summary Get colour variants under model
// @Description Get colour variants under model
// @Tags Admin/Product_Management/Colour_Variant
// @Tags Seller/Product_Management/Colour_Variant
// @Tags User/Browse
// @Produce json
// @Security BearerTokenAuth
// @Param modelID path int true "Model ID"
// @Success 200 {object} response.GetColourVariantsUnderModelResponse
// @Failure 400 {object} response.SME{}
// @Router /admin/colourvariants/{modelID} [get]
func (cvh *ProductHandler) GetColourVariantsUnderModel(c *gin.Context) {

	modelIDParam := c.Param("modelID")
	modelIDstring, errr := strconv.ParseUint(modelIDParam, 10, 64)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error parsing modelID:", errr))
		return
	}

	modelID := uint(modelIDstring)
	colourVariants, err := cvh.productUseCase.GetColourVariantsUnderModel(modelID)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.GetColourVariantsUnderModelResponse{
		ColourVariants: *colourVariants,
	})
}
