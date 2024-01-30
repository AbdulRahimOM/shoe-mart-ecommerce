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
	
	var req requestModels.AddModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
	}

	//add model
	if err := h.modelsUseCase.AddModel(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding model. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSME("Model added successfully"))
}

func (h *ModelHandler) GetModelsByBrandsAndCategories(c *gin.Context) {

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
			c.JSON(http.StatusBadRequest, response.FailedSME("Invalid request. Try Again", err))
			return
		}
	}

	if categoryIDParam != "" {
		categoryExists = true
		categoryIDs, err = requestValidation.ValidateAndParseIDs(categoryIDParam)
		if err != nil {
			fmt.Println("error parsing category id. error:", err)
			c.JSON(http.StatusBadRequest, response.FailedSME("Invalid request. Try Again", err))
			return
		}
	}

	fmt.Println("brandIDs:", brandIDs)
	fmt.Println("categoryIDs:", categoryIDs)

	//get models
	var models *[]entities.Models
	models, err = h.modelsUseCase.GetModelsByBrandsAndCategories(brandExists, brandIDs, categoryExists, categoryIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error fetching models. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.GetModelsResponse{
		Status:  "success",
		Message: "Models fetched successfully",
		Models: *models,
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
	
	var req requestModels.EditModelReq
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

	//edit model name
	if err := h.modelsUseCase.EditModelName(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error editing model name. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSME("Model name edited successfully"))
}
