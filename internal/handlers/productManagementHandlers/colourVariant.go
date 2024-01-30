package productManagementHandlers

import (
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// add colour variant handler
// @Summary Add colour variant
// @Description Add colour variant
// @Tags admin
// @Accept json
// @Produce json
// @Param addColourVariantReq body requestModels.AddColourVariantReq true "Add Colour Variant Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/addcolourvariant [post]
func (cvh *ProductHandler) AddColourVariant(c *gin.Context) {

	var req requestModels.AddColourVariantReq
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
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

	c.JSON(http.StatusOK, response.SuccessSME("Colour variant added successfully"))
}

// edit colour variant handler
func (cvh *ProductHandler) EditColourVariant(c *gin.Context) {

	var req requestModels.EditColourVariantReq
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
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//add colour variant
	if err := cvh.productUseCase.EditColourVariant(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error editing colour variant. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSME("Colour variant edited successfully"))
}

// get colour variants under model handler
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
