package reporthandlers

import (
	"net/http"
	"strconv"

	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"

	"github.com/gin-gonic/gin"
)

// TopProductsHandler
// @Summary Get top products
// @Description Get top products
// @Tags Admin/Analytics/Top_Performers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param l query int false "Limit"
// @Success 200 {object} response.TopProductsResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/top-products [get]
func (h *ReportsHandler) TopProductsHandler(c *gin.Context) {

	//get query param "l"
	limit64, errr := strconv.ParseUint(c.Query("l"), 10, 64)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.FromError(errr))
		return
	}
	limit := int(limit64)

	topProducts, err := h.reportsUC.GetTopProducts(limit)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopProductsResponse{
			Limit:       limit,
			TopProducts: *topProducts,
		})
	}
}

// TopSellersHandler
// @Summary Get top sellers
// @Description Get top sellers
// @Tags Admin/Analytics/Top_Performers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param l query int false "Limit"
// @Success 200 {object} response.TopSellersResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/top-sellers [get]
func (h *ReportsHandler) TopSellersHandler(c *gin.Context) {

	//get query param "l"
	limit64, errr := strconv.ParseUint(c.Query("l"), 10, 64)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.FromError(errr))
		return
	}
	limit := int(limit64)

	topSellers, err := h.reportsUC.GetTopSellers(limit)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopSellersResponse{
			Limit:      limit,
			TopSellers: *topSellers,
		})
	}
}

// TopBrandsHandler
// @Summary Get top brands
// @Description Get top brands
// @Tags Admin/Analytics/Top_Performers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param l query int false "Limit"
// @Success 200 {object} response.TopBrandsResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/top-brands [get]
func (h *ReportsHandler) TopBrandsHandler(c *gin.Context) {

	//get query param "l"
	limit64, errr := strconv.ParseUint(c.Query("l"), 10, 64)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.FromError(errr))
		return
	}
	limit := int(limit64)

	topBrands, err := h.reportsUC.GetTopBrands(limit)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopBrandsResponse{
			Limit:     limit,
			TopBrands: *topBrands,
		})
	}
}

// TopModelsHandler
// @Summary Get top models
// @Description Get top models
// @Tags Admin/Analytics/Top_Performers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param l query int false "Limit"
// @Success 200 {object} response.TopModelsResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/top-models [get]
func (h *ReportsHandler) TopModelsHandler(c *gin.Context) {

	//get query param "l"
	limit64, errr := strconv.ParseUint(c.Query("l"), 10, 64)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.FromError(errr))
		return
	}
	limit := int(limit64)

	topModels, err := h.reportsUC.GetTopModels(limit)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopModelsResponse{
			Limit:     limit,
			TopModels: *topModels,
		})
	}
}
