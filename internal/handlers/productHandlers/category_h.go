package producthandler

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecase "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUseCase usecase.ICategoryUC
}

func NewCategoryHandler(usecase usecase.ICategoryUC) *CategoryHandler {
	return &CategoryHandler{CategoryUseCase: usecase}
}

// add category handler
// @Summary Add category
// @Description Add category
// @Tags Admin/Product_Management/Category
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addCategoryReq body req.AddCategoryReq{} true "Add Category Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/addcategory [post]
func (h *CategoryHandler) AddCategory(c *gin.Context) {

	var req request.AddCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	err := h.CategoryUseCase.AddCategory(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, nil)
		return
	}
}

// get categories handler
// @Summary Get categories
// @Description Get categories
// @Tags Admin/Product_Management/Category
// @Tags Seller/Product_Management/Category
// @Tags User/Browse
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetCategoriesResponse
// @Failure 400 {object} response.SME{}
// @Router /admin/categories [get]
// @Router /seller/categories [get]
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {

	var categories *[]entities.Categories
	categories, err := h.CategoryUseCase.GetCategories()
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		fmt.Println("categories fetched successfully")
		c.JSON(http.StatusOK, response.GetCategoriesResponse{
			Categories: *categories,
		})
		return
	}
}

// edit category handler
// @Summary Edit category
// @Description Edit category
// @Tags Admin/Product_Management/Category
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param editCategoryReq body req.EditCategoryReq{} true "Edit Category Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/editcategory [patch]
func (h *CategoryHandler) EditCategory(c *gin.Context) {

	var req request.EditCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	err := h.CategoryUseCase.EditCategory(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK,  nil)
		return
	}
}
