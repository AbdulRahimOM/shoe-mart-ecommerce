package producthandler

import (
	e "MyShoo/internal/domain/customErrors"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//get sellerID from token
	sellerID, err := tools.GetSellerID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error occurred. Try Again", err))
		return
	}

	//image upload handling
	formFile, err := c.FormFile("imageUrl")
	if err != nil {
		fmt.Println("error getting image file from request. err: ", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error uploading image file. Try Again", err))
		return
	}

	path := filepath.Join(os.TempDir(), formFile.Filename)
	if err := c.SaveUploadedFile(formFile, path); err != nil {
		fmt.Println("error saving image file to temp dir. err: ", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error uploading image file. Try Again", err))
		return
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening image file. err: ", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error uploading image file. Try Again", err))
		return
	}
	defer file.Close()

	//add colour variant
	if err := cvh.productUseCase.AddColourVariant(sellerID, &req, file); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding colour variant. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Colour variant added successfully"))
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
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//add colour variant
	if err := cvh.productUseCase.EditColourVariant(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error editing colour variant. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Colour variant edited successfully"))
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
	modelIDstring, err := strconv.ParseUint(modelIDParam, 10, 64)
	if err != nil {
		fmt.Println("error parsing modelID. err: ", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error occured. Try Again", err))
		return
	}

	modelID := uint(modelIDstring)
	colourVariants, err := cvh.productUseCase.GetColourVariantsUnderModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error fetching colour variants under model. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.GetColourVariantsUnderModelResponse{
		Status:         "success",
		Message:        "Colour variants under model fetched successfully",
		ColourVariants: *colourVariants,
	})
}
