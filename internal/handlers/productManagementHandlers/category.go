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
	fmt.Println("==================\nentered add category handler")
	var addCategoryReq requestModels.AddCategoryReq
	if err := c.ShouldBindJSON(&addCategoryReq); err != nil {
		fmt.Println("error binding request")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(addCategoryReq); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Provide valid data",
			Error:   errResponse,
		})
		return
	}

	err := h.CategoryUseCase.AddCategory(&addCategoryReq)
	if err != nil {
		fmt.Println("error from usecase")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Couldn't add category. Some error occured. Try Again",
			Error:   err.Error(),
		})
		return
	} else {
		fmt.Println("category added successfully")
		c.JSON(http.StatusOK, response.SME{
			Status:  "success",
			Message: "Category added successfully",
			Error:   "",
		})
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
	fmt.Println("==================\nentered get categories handler")
	var categories *[]entities.Categories
	categories, err := h.CategoryUseCase.GetCategories()
	if err != nil {
		fmt.Println("error from usecase")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Couldn't get categories. Some error occured. Try Again",
			Error:   err.Error(),
		})
		return
	} else {
		fmt.Println("categories fetched successfully")
		c.JSON(http.StatusOK, response.GetCategoriesResponse{
			Status:     "success",
			Message:    "Categories fetched successfully",
			Categories: *categories,
			Error:      "",
		})
		return
	}
}

//edit category handler
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
	fmt.Println("==================\nentered edit category handler")
	var req requestModels.EditCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("error binding request")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Provide valid data",
			Error:   errResponse,
		})
		return
	}

	err := h.CategoryUseCase.EditCategory(&req)
	if err != nil {
		fmt.Println("error from usecase")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Couldn't edit category. Some error occured. Try Again",
			Error:   err.Error(),
		})
		return
	} else {
		fmt.Println("category edited successfully")
		c.JSON(http.StatusOK, response.SME{
			Status:  "success",
			Message: "Category edited successfully",
			Error:   "",
		})
		return
	}
}