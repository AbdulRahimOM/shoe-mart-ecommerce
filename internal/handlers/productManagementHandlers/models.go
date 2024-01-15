package productManagementHandlers

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	modelsUseCase usecaseInterface.IModelsUC
}

func NewModelHandler(modelsUseCase usecaseInterface.IModelsUC) *ModelHandler {
	return &ModelHandler{modelsUseCase: modelsUseCase}
}

// models handler
// @Summary Add model
// @Description Add model
// @Tags admin
// @Accept json
// @Produce json
// @Param addModelReq body requestModels.AddModelReq true "Add Model Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/addmodel [post]
func (h *ModelHandler) AddModel(c *gin.Context) {
	fmt.Println("Handler ::: add model handler")
	var req requestModels.AddModelReq

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

	//add model
	if err := h.modelsUseCase.AddModel(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding model. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Model added successfully",
	})
}

func (h *ModelHandler) GetModelsByBrandsAndCategories(c *gin.Context) {
	fmt.Println("Handler ::: GetModelsByBrandsAndCategories handler")

	// Get parameters from URL
	brandIDParam := c.Query("brandID")
	categoryIDParam := c.Query("categoryID")

	var brandExists, categoryExists bool
	var brandIDs, categoryIDs []uint
	var err error

	// Validate and convert the string parameters to arrays of integers
	if brandIDParam != "" {
		brandExists = true
		brandIDs, err = requestValidation.ValidateAndParseIDs(brandIDParam)
		if err != nil {
			fmt.Println("error parsing brand id. error:", err)
			c.JSON(http.StatusBadRequest, response.GetModelsResponse{
				Status:  "failed",
				Message: "Error parsing brand id. Try Again",
				Error:   err.Error(),
			})
			return
		}
	}

	if categoryIDParam != "" {
		categoryExists = true
		categoryIDs, err = requestValidation.ValidateAndParseIDs(categoryIDParam)
		if err != nil {
			fmt.Println("error parsing category id. error:", err)
			c.JSON(http.StatusBadRequest, response.GetModelsResponse{
				Status:  "failed",
				Message: "Error parsing category id. Try Again",
				Error:   err.Error(),
			})
			return
		}
	}

	fmt.Println("brandIDs:", brandIDs)
	fmt.Println("categoryIDs:", categoryIDs)
	//get models
	var models *[]entities.Models
	fmt.Println("brandExists:", brandExists, "categoryExists:", categoryExists)
	models, err = h.modelsUseCase.GetModelsByBrandsAndCategories(brandExists, brandIDs, categoryExists, categoryIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting models. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetModelsResponse{
		Status:  "success",
		Message: "Models fetched successfully",
		Error:   "",
		Models:  *models,
	})

}

// edit model name handler
// @Summary Edit model
// @Description Edit model
// @Tags admin
// @Accept json
// @Produce json
// @Param editModelNameReq body requestModels.EditModelNameReq true "Edit Model Name Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/editmodelname [patch]
func (h *ModelHandler) EditModel(c *gin.Context) {
	fmt.Println("Handler ::: edit model name handler")
	var req requestModels.EditModelReq

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

	//edit model name
	if err := h.modelsUseCase.EditModelName(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing model name. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Model name edited successfully",
	})
}
