package productManagementHandlers

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase usecaseInterface.IProductsUC
}

func NewProductHandler(productUseCase usecaseInterface.IProductsUC) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

// get products handler
// @Summary Get products
// @Description Get products
// @Tags Admin/Product_Management/Products
// @Tags Seller/Product_Management/Products
// @Tags User/Browse
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.SMED{}
// @Failure 400 {object} response.SME{}
// @Router /seller/products [get]
// @Router /admin/products [get]
// @Router /user/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {

	//get products
	var products *[]response.ResponseProduct
	products, err := h.productUseCase.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error fetching products. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SMED{
		Status:  "success",
		Message: "Products fetched successfully",
		Data:    products,
	})

}

// add stock handler
// @Summary Add stock
// @Description Add stock
// @Tags Seller/Product_Management/Stock
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addStockReq body requestModels.AddStockReq{} true "Add Stock Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /seller/addstock [post]
func (h *ProductHandler) AddStock(c *gin.Context) {

	var req requestModels.AddStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//check if sellerID in token and request body match
	sellerID, err := tools.GetSellerID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error occured", err))
		return
	}
	if sellerID != req.SellerID {
		fmt.Println("Seller ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "Seller ID in token and request body do not match",
		})
		return
	}

	//add stock
	if err := h.productUseCase.AddStock(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding stock. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Stock added successfully"))
}

// edit stock handler
// @Summary Edit stock
// @Description Edit stock
// @Tags Seller/Product_Management/Stock
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param editStockReq body requestModels.EditStockReq{} true "Edit Stock Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /seller/editstock [patch]
func (h *ProductHandler) EditStock(c *gin.Context) {

	var req requestModels.EditStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//check if sellerID in token and request body match
	sellerID, err := tools.GetSellerID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing stock. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if sellerID != req.SellerID {
		fmt.Println("Seller ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "Seller ID in token and request body do not match",
		})
		return
	}

	//add stock
	if err := h.productUseCase.EditStock(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error editing stock. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Stock edited successfully"))
}
