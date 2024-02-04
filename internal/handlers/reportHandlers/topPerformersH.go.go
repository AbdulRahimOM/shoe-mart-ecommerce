package reporthandlers

import (
	response "MyShoo/internal/models/responseModels"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *ReportsHandler) TopProductsHandler(c *gin.Context) {

	//get query param "l"
	limit64, err := strconv.ParseUint(c.Query("l"), 10, 64)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting limit from query", err))
		return
	}
	limit := int(limit64)

	topProducts, err := h.reportsUC.GetTopProducts(limit)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting top products", err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopProductsResponse{
			Status:      "success",
			Limit:       limit,
			TopProducts: *topProducts,
		})
	}
}

func (h *ReportsHandler) TopSellersHandler(c *gin.Context) {

	//get query param "l"
	limit64, err := strconv.ParseUint(c.Query("l"), 10, 64)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting limit from query", err))
		return
	}
	limit := int(limit64)

	topSellers, err := h.reportsUC.GetTopSellers(limit)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting top sellers", err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopSellersResponse{
			Status:     "success",
			Limit:      limit,
			TopSellers: *topSellers,
		})
	}
}

func (h *ReportsHandler) TopBrandsHandler(c *gin.Context) {

	//get query param "l"
	limit64, err := strconv.ParseUint(c.Query("l"), 10, 64)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting limit from query", err))
		return
	}
	limit := int(limit64)

	topBrands, err := h.reportsUC.GetTopBrands(limit)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting top brands", err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopBrandsResponse{
			Status:    "success",
			Limit:     limit,
			TopBrands: *topBrands,
		})
	}
}

func (h *ReportsHandler) TopModelsHandler(c *gin.Context) {

	//get query param "l"
	limit64, err := strconv.ParseUint(c.Query("l"), 10, 64)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting limit from query", err))
		return
	}
	limit := int(limit64)

	topModels, err := h.reportsUC.GetTopModels(limit)
	if err != nil {
		c.JSON(400, response.FailedSME("error occured in getting top models", err))
		return
	} else {
		c.JSON(http.StatusOK, response.TopModelsResponse{
			Status:    "success",
			Limit:     limit,
			TopModels: *topModels,
		})
	}
}
