package producthandler

import (
	"net/http"

	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
	requestValidation "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/validation"

	"github.com/gin-gonic/gin"
)

type BrandsHandler struct {
	brandsUseCase usecase.IBrandsUC
}

func NewBrandHandler(uc usecase.IBrandsUC) *BrandsHandler {
	return &BrandsHandler{brandsUseCase: uc}
}

// add brands handler
// @Summary Add brand
// @Description Add brand
// @Tags Seller/Product_Management/Brand
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addBrandReq body req.AddBrandReq{} true "Add Brand Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /seller/addbrand [post]
func (bh *BrandsHandler) AddBrand(c *gin.Context) {

	var req request.AddBrandReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//add brand
	if err := bh.brandsUseCase.AddBrand(&req); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// get brands handler
// @Summary Get brands
// @Description Get brands (categorised by first alphabet)
// @Tags Admin/Product_Management/Brand
// @Tags Seller/Product_Management/Brand
// @Tags User/Browse
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetBrandsResponse
// @Failure 400 {object} response.SME{}
// @Router /admin/brands [get]
// @Router /seller/brands [get]
// @Router /brands [get]
func (bh *BrandsHandler) GetBrands(c *gin.Context) {

	var brands *[26]entities.BrandsByAlphabet
	brands, err := bh.brandsUseCase.GetBrands()
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.GetBrandsResponse{
		BrandsByAlphabet: *brands,
	})
}

// edit brand handler
// @Summary Edit brand
// @Description Edit brand
// @Tags Admin/Product_Management/Brand
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param editBrandReq body req.EditBrandReq{} true "Edit Brand Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/editbrand [patch]
func (bh *BrandsHandler) EditBrand(c *gin.Context) {

	var req request.EditBrandReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//edit brand
	if err := bh.brandsUseCase.EditBrand(&req); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("brand edited"))
}
