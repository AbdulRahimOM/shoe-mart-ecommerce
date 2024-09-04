package orderhandler

import (
	"net/http"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	requestValidation "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/validation"

	"github.com/gin-gonic/gin"
)

// Admin can add new coupons
// @Summary Add new coupon
// @Description Admin can add new coupon
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param newCouponReq body req.NewCouponReq{} true "New Coupon Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/newcoupon [post]
func (h *OrderHandler) NewCouponHandler(c *gin.Context) {

	//get req from body
	var req request.NewCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//call usecase
	err := h.orderUseCase.CreateNewCoupon(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.MsgAndError("error adding coupon", err))
	} else {
		c.JSON(200, response.SuccessSM("coupon added"))
	}
}

// BlockCouponHandler
// @Summary Block coupon
// @Description Admin can block(suspend) a coupon
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param blockCouponReq body req.BlockCouponReq{} true "Block Coupon Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/blockcoupon [patch]
func (h *OrderHandler) BlockCouponHandler(c *gin.Context) {

	//get req from body
	var req request.BlockCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//call usecase
	err := h.orderUseCase.BlockCoupon(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.MsgAndError("error blocking coupon", err))
	} else {
		c.JSON(200, response.SuccessSM("coupon blocked"))
	}
}

// UnblockCouponHandler
// @Summary Unblock coupon
// @Description Admin can unblock(re-activate) a coupon
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param unblockCouponReq body req.UnblockCouponReq{} true "Unblock Coupon Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/unblockcoupon [patch]
func (h *OrderHandler) UnblockCouponHandler(c *gin.Context) {

	//get req from body
	var req request.UnblockCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//call usecase
	err := h.orderUseCase.UnblockCoupon(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.MsgAndError("error unblocking coupon", err))
	} else {
		c.JSON(200, response.SuccessSM("coupon unblocked"))
	}
}

// GetCoupons
// @Summary Get coupons
// @Description Admin can get all coupons, active coupons, expired coupons, upcoming coupons
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param criteria query string true "all, active, expired, upcoming"
// @Success 200 {object} response.GetCouponRes{}
// @Failure 400 {object} response.SME{}
// @Router /admin/coupons [get]
func (h *OrderHandler) GetCoupons(c *gin.Context) {

	criteria := c.Query("criteria")
	var coupons *[]entities.Coupon
	var err *e.Error
	switch criteria {
	case "all":
		coupons, err = h.orderUseCase.GetAllCoupons()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		}
	case "expired":
		coupons, err = h.orderUseCase.GetExpiredCoupons()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		}
	case "active":
		coupons, err = h.orderUseCase.GetActiveCoupons()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		}
	case "upcoming":
		coupons, err = h.orderUseCase.GetUpcomingCoupons()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		}
	default:
		c.JSON(400, response.FailedSME("invalid url parameter", nil))
	}

	c.JSON(200, response.GetCouponRes{
		Coupons: *coupons,
	})

}
