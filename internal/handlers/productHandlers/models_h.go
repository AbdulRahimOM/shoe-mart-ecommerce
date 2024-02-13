package producthandler

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecase "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	modelsUseCase usecase.IModelsUC
}

func NewModelHandler(modelsUseCase usecase.IModelsUC) *ModelHandler {
	return &ModelHandler{modelsUseCase: modelsUseCase}
}

// models handler
// @Summary Add model
// @Description Add model
// @Tags Seller/Product_Management/Model
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addModelReq body req.AddModelReq true "Add Model Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /seller/addmodel [post]
func (h *ModelHandler) AddModel(c *gin.Context) {

	var req request.AddModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//add model
	if err := h.modelsUseCase.AddModel(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding model. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Model added successfully"))
}

// get models by brands and categories handler
// @Summary Get models by brands and categories
// @Description Get models by brands and categories
// @Tags Admin/Product_Management/Model
// @Tags Seller/Product_Management/Model
// @Tags User/Browse
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param brandID query string false "Brand ID"
// @Param categoryID query string false "Category ID"
// @Success 200 {object} response.GetModelsResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/models [get]
// @Router /seller/models [get]
// @Router /models [get]
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
		Models:  *models,
	})

}

// edit model name handler
// @Summary Edit model
// @Description Edit model
// @Tags Admin/Product_Management/Model
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param editModelReq body req.EditModelReq{} true "Edit Model Name Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/editmodel [patch]
func (h *ModelHandler) EditModel(c *gin.Context) {

	var req request.EditModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//edit model name
	if err := h.modelsUseCase.EditModelName(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error editing model name. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Model name edited successfully"))
}
