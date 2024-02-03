package ordermanagementHandlers

import (
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) NewCouponHandler(c *gin.Context) {

	//get req from body
	var newCouponReq requestModels.NewCouponReq
	if err := c.ShouldBindJSON(&newCouponReq); err != nil {
		c.JSON(400, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(newCouponReq); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(400, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//call usecase
	message, err := h.orderUseCase.CreateNewCoupon(&newCouponReq)
	if err != nil {
		c.JSON(400, response.FailedSME(message, err))
		return
	}

	c.JSON(200, response.SuccessSME(message))
}

// BlockCouponHandler
func (h *OrderHandler) BlockCouponHandler(c *gin.Context) {

	//get req from body
	var blockCouponReq requestModels.BlockCouponReq
	if err := c.ShouldBindJSON(&blockCouponReq); err != nil {
		c.JSON(400, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(blockCouponReq); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(400, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//call usecase
	message, err := h.orderUseCase.BlockCoupon(&blockCouponReq)
	if err != nil {
		c.JSON(400, response.FailedSME(message, err))
		return
	}

	c.JSON(200, response.SuccessSME(message))
}

// UnblockCouponHandler
func (h *OrderHandler) UnblockCouponHandler(c *gin.Context) {

	//get req from body
	var unblockCouponReq requestModels.UnblockCouponReq
	if err := c.ShouldBindJSON(&unblockCouponReq); err != nil {
		c.JSON(400, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(unblockCouponReq); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(400, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//call usecase
	message, err := h.orderUseCase.UnblockCoupon(&unblockCouponReq)
	if err != nil {
		c.JSON(400, response.FailedSME(message, err))
		return
	}

	c.JSON(200, response.SuccessSME(message))
}

// GetCoupons
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
