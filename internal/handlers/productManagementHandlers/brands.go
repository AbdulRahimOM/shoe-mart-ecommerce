package productManagementHandlers

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BrandsHandler struct {
	brandsUseCase usecaseInterface.IBrandsUC
}

func NewBrandHandler(uc usecaseInterface.IBrandsUC) *BrandsHandler {
	return &BrandsHandler{brandsUseCase: uc}
}

// add brands handler
// @Summary Add brand
// @Description Add brand
// @Tags admin
// @Accept json
// @Produce json
// @Param addBrandReq body requestModels.AddBrandReq true "Add Brand Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/addbrand [post]
func (bh *BrandsHandler) AddBrand(c *gin.Context) {

	var req requestModels.AddBrandReq
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

	//add brand
	if err := bh.brandsUseCase.AddBrand(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding brand. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Brand added successfully"))
}

// get brands handler
// @Summary Get brands
// @Description Get brands
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/getbrands [get]
func (bh *BrandsHandler) GetBrands(c *gin.Context) {

	var brands *[26]entities.BrandsByAlphabet
	brands, err := bh.brandsUseCase.GetBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error fetching brands. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.GetBrandsResponse{
		Status:           "success",
		Message:          "Brands fetched successfully",
		BrandsByAlphabet: *brands,
	})
}

// edit brand handler
// @Summary Edit brand
// @Description Edit brand
// @Tags admin
// @Accept json
// @Produce json
// @Param editBrandReq body requestModels.EditBrandReq true "Edit Brand Request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/editbrand [post]
func (bh *BrandsHandler) EditBrand(c *gin.Context) {

	var req requestModels.EditBrandReq
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

	//edit brand
	if err := bh.brandsUseCase.EditBrand(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error editing brand. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Brand edited successfully"))
}
