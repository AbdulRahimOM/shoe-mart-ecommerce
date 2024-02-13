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

type CategoryHandler struct {
	CategoryUseCase usecaseInterface.ICategoryUC
}

func NewCategoryHandler(usecase usecaseInterface.ICategoryUC) *CategoryHandler {
	return &CategoryHandler{CategoryUseCase: usecase}
}

// add category handler
// @Summary Add category
// @Description Add category
// @Tags admin
// @Accept json
// @Produce json
// @Param addCategoryReq body requestModels.AddCategoryReq true "Add Category Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/addcategory [post]
func (h *CategoryHandler) AddCategory(c *gin.Context) {

	var addCategoryReq requestModels.AddCategoryReq
	if err := c.ShouldBindJSON(&addCategoryReq); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(addCategoryReq); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Provide valid data", errResponse))
		return
	}

	err := h.CategoryUseCase.AddCategory(&addCategoryReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Couldn't add category. Try Again", err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("Category added successfully"))
		return
	}
}

// get categories handler
// @Summary Get categories
// @Description Get categories
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/getcategories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {

	var categories *[]entities.Categories
	categories, err := h.CategoryUseCase.GetCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Couldn't fetch categories. Try Again", err))
		return
	} else {
		fmt.Println("categories fetched successfully")
		c.JSON(http.StatusOK, response.GetCategoriesResponse{
			Status:     "success",
			Message:    "Categories fetched successfully",
			Categories: *categories,
		})
		return
	}
}

// edit category handler
// @Summary Edit category
// @Description Edit category
// @Tags admin
// @Accept json
// @Produce json
// @Param editCategoryReq body requestModels.EditCategoryReq true "Edit Category Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/editcategory [post]
func (h *CategoryHandler) EditCategory(c *gin.Context) {

	var req requestModels.EditCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Provide valid data", errResponse))
		return
	}

	err := h.CategoryUseCase.EditCategory(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Couldn't edit category. Try Again", err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("Category edited successfully"))
		return
	}
}
