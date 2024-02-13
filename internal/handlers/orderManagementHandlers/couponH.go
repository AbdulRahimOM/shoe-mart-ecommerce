package ordermanagementHandlers

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin can add new coupons
// @Summary Add new coupon
// @Description Admin can add new coupon
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param newCouponReq body requestModels.NewCouponReq{} true "New Coupon Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/newcoupon [post]
func (h *OrderHandler) NewCouponHandler(c *gin.Context) {

	//get req from body
	var req requestModels.NewCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//call usecase
	message, err := h.orderUseCase.CreateNewCoupon(&req)
	if err != nil {
		c.JSON(400, response.FailedSME(message, err))
		return
	}

	c.JSON(200, response.SuccessSM(message))
}

// BlockCouponHandler
// @Summary Block coupon
// @Description Admin can block(suspend) a coupon
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param blockCouponReq body requestModels.BlockCouponReq{} true "Block Coupon Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/blockcoupon [patch]
func (h *OrderHandler) BlockCouponHandler(c *gin.Context) {

	//get req from body
	var req requestModels.BlockCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//call usecase
	message, err := h.orderUseCase.BlockCoupon(&req)
	if err != nil {
		c.JSON(400, response.FailedSME(message, err))
		return
	}

	c.JSON(200, response.SuccessSM(message))
}

// UnblockCouponHandler
// @Summary Unblock coupon
// @Description Admin can unblock(re-activate) a coupon
// @Tags Admin/Coupon
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param unblockCouponReq body requestModels.UnblockCouponReq{} true "Unblock Coupon Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/unblockcoupon [patch]
func (h *OrderHandler) UnblockCouponHandler(c *gin.Context) {

	//get req from body
	var req requestModels.UnblockCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//call usecase
	message, err := h.orderUseCase.UnblockCoupon(&req)
	if err != nil {
		c.JSON(400, response.FailedSME(message, err))
		return
	}

	c.JSON(200, response.SuccessSM(message))
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
	var message string
	var err error
	switch criteria {
	case "all":
		coupons, message, err = h.orderUseCase.GetAllCoupons()
		if err != nil {
			c.JSON(400, response.FailedSME(message, err))
			return
		}

	case "expired":
		coupons, message, err = h.orderUseCase.GetExpiredCoupons()
		if err != nil {
			c.JSON(400, response.FailedSME(message, err))
			return
		}
	case "active":
		coupons, message, err = h.orderUseCase.GetActiveCoupons()
		if err != nil {
			c.JSON(400, response.FailedSME(message, err))
			return
		}
	case "upcoming":
		coupons, message, err = h.orderUseCase.GetUpcomingCoupons()
		if err != nil {
			c.JSON(400, response.FailedSME(message, err))
		}
	default:
		c.JSON(400, response.FailedSME(msg.InvalidRequest, errors.New("invalid url parameter")))
	}

	c.JSON(200, response.GetCouponRes{
		Status:  "success",
		Message: "Coupons fetched successfully",
		Coupons: *coupons,
	})

}
