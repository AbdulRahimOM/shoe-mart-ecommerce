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
	fmt.Println("Handler ::: add colour variant handler")

	var req requestModels.AddColourVariantReq
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		fmt.Println("error binding request. err: ", err)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}
	fmt.Println("req: ", req)

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

	//get sellerID from token
	sellerID, err := tools.GetSellerID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error occured. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//image upload handling

	formFile, err := c.FormFile("imageUrl")
	if err != nil {
		fmt.Println("error getting image file from request. err: ", err)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error uploading image file. Try Again",
			Error:   err.Error(),
		})
		return
	}
	// fmt.Println("formFile: ", formFile)

	path := filepath.Join(os.TempDir(), formFile.Filename)
	if err := c.SaveUploadedFile(formFile, path); err != nil {
		fmt.Println("error saving image file to temp dir. err: ", err)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error uploading image file. Try Again",
			Error:   err.Error(),
		})
		return
	}
	// fmt.Println("path: ", path)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening image file. err: ", err)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error uploading image file. Try Again",
			Error:   err.Error(),
		})
		return
	}
	fmt.Println("H file: ", file)
	defer file.Close()

	//add colour variant
	if err := cvh.productUseCase.AddColourVariant(sellerID, &req, file); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding colour variant. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Colour variant added successfully",
	})
}

// edit colour variant handler
func (cvh *ProductHandler) EditColourVariant(c *gin.Context) {
	fmt.Println("Handler ::: add colour variant handler")
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
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//add colour variant
	if err := cvh.productUseCase.EditColourVariant(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding colour variant. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Colour variant added successfully",
	})
}

// get colour variants under model handler
func (cvh *ProductHandler) GetColourVariantsUnderModel(c *gin.Context) {
	fmt.Println("Handler ::: get colour variants under model handler")

	modelIDParam := c.Param("modelID")
	modelIDstring, err := strconv.ParseUint(modelIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error parsing modelID. Try Again",
			Error:   err.Error(),
		})
		return
	}

	modelID := uint(modelIDstring)
	// var req requestModels.GetColourVariantsUnderModelReq

	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, response.SME{
	// 		Status:  "failed",
	// 		Message: "Error binding request. Try Again",
	// 		Error:   err.Error(),
	// 	})
	// 	return
	// }

	//validate request
	// if err := requestValidation.ValidateRequest(req); err != nil {
	// 	errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
	// 	fmt.Println(errResponse)
	// 	c.JSON(http.StatusBadRequest, response.SME{
	// 		Status:  "failed",
	// 		Message: "Error validating request. Try Again",
	// 		Error:   errResponse,
	// 	})
	// 	return
	// }

	colourVariants, err := cvh.productUseCase.GetColourVariantsUnderModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting colour variants under model. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetColourVariantsUnderModelResponse{
		Status:         "success",
		Message:        "Colour variants under model fetched successfully",
		ColourVariants: *colourVariants,
	})
}
