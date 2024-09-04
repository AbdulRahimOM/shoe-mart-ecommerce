package producthandler

import (
	"net/http"

	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/tools"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
	requestValidation "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/validation"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase usecase.IProductsUC
}

func NewProductHandler(productUseCase usecase.IProductsUC) *ProductHandler {
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
	products, err := h.productUseCase.GetProducts()
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SMED{
		Data: products,
	})

}

// add stock handler
// @Summary Add stock
// @Description Add stock
// @Tags Seller/Product_Management/Stock
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addStockReq body req.AddStockReq{} true "Add Stock Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /seller/addstock [post]
func (h *ProductHandler) AddStock(c *gin.Context) {

	var req request.AddStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//check if sellerID in token and request body match
	sellerID, errr := tools.GetSellerID(c)
	if errr != nil {
		c.JSON(http.StatusForbidden, response.MsgAndError("error getting sellerID from token:", errr))
		return
	}
	if sellerID != req.SellerID {
		// fmt.Println("Seller ID in token and request body do not match. Corrupted request!!")
		// c.JSON(http.StatusBadRequest, response.SME{
		// 	Status:  "failed",
		// 	Message: "Corrupted request. Try Again",
		// 	Error:   "Seller ID in token and request body do not match",
		// })
		c.JSON(http.StatusForbidden, response.FromErrByText("seller ID in token and request body do not match"))
		return
	}

	//add stock
	if err := h.productUseCase.AddStock(&req); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("stock added"))
}

// edit stock handler
// @Summary Edit stock
// @Description Edit stock
// @Tags Seller/Product_Management/Stock
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param editStockReq body req.EditStockReq{} true "Edit Stock Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /seller/editstock [patch]
func (h *ProductHandler) EditStock(c *gin.Context) {

	var req request.EditStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//check if sellerID in token and request body match
	sellerID, errr := tools.GetSellerID(c)
	if errr != nil {
		c.JSON(http.StatusForbidden, response.MsgAndError("error getting sellerID from token:", errr))
		return
	}
	if sellerID != req.SellerID {
		c.JSON(http.StatusForbidden, response.FromErrByText("seller ID in token and request body do not match"))
		return
	}

	//add stock
	if err := h.productUseCase.EditStock(&req); err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("stock edited"))
}
