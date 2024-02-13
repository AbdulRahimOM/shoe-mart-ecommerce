package ordermanagementHandlers

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"strconv"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase usecaseInterface.IOrderUC
}

func NewOrderHandler(orderUseCase usecaseInterface.IOrderUC) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// MakeOrder
// @Summary Make Order
// @Description Make Order
// @Tags User/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param makeOrderReq body requestModels.MakeOrderReq{} true "Make Order Request"
// @Success 201 {object} response.CODOrderResponse
// @Success 201 {object} response.OnlinePaymentOrderResponse
// @Failure 400 {object} response.SME{}
// @Router /makeorder [post]
func (h *OrderHandler) MakeOrder(c *gin.Context) {

	var req *requestModels.MakeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error making order. Try Again", err))
		return
	}
	req.UserID = userID

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//make order
	orderInfo, proceedToPaymentInfo, message, err := h.orderUseCase.MakeOrder(req)
	if err != nil {
		switch err.Error() {
		case "cart is empty":
			c.JSON(http.StatusBadRequest, response.FailedSME(message, err))
		case "stock not available":
			c.JSON(http.StatusForbidden, response.FailedSME(message, err))
		default:
			c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		}
		return

	}
	if proceedToPaymentInfo == nil {
		c.JSON(http.StatusCreated, response.CODOrderResponse{
			Status:    "success",
			Message:   message,
			OrderInfo: *orderInfo,
		})
	} else {
		c.JSON(http.StatusCreated, response.OnlinePaymentOrderResponse{
			Status:               "success",
			Message:              message,
			OrderInfo:            *orderInfo,
			ProceedToPaymentInfo: *proceedToPaymentInfo,
		})
	}
}

// Get Orders of the user
// @Summary Get Orders of the user
// @Description Get Orders of the user
// @Tags User/Order
// @Produce json
// @Security BearerTokenAuth
// @Param p query string false "page number"
// @Param l query string false "limit"
// @Success 200 {object} response.GetOrdersResponse
// @Failure 400 {object} response.SME{}
// @Router /myorders [get]
func (h *OrderHandler) GetOrdersOfUser(c *gin.Context) {

	//get pagination params
	page := c.Query("p")
	limit := c.Query("l")
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	// Validate and convert the string parameters to integers
	pageInt, err1 := requestValidation.ValidateAndParseInt(page)
	limitInt, err2 := requestValidation.ValidateAndParseInt(limit)
	if err1 != nil || err2 != nil {
		fmt.Println("error parsing p parameter. error:", err1, ",", err2)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error in url. Try Again", fmt.Errorf("%v, %v", err1, err2)))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting orders. Try Again", err))
		return
	}

	//get orders
	var orders *[]response.ResponseOrderInfo
	orders, message, err := h.orderUseCase.GetOrdersOfUser(userID, pageInt, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.GetOrdersResponse{
		Status:  "success",
		Message: "Orders fetched successfully",

		OrdersInfo: *orders,
	})
}

// Get All Orders (for admin)
// @Summary Get All Orders
// @Description Get All Orders (for admin)
// @Tags Admin/Order
// @Produce json
// @Security BearerTokenAuth
// @Param p query string false "page number"
// @Param l query string false "limit"
// @Success 200 {object} response.GetOrdersResponse
// @Failure 400 {object} response.SME{}
// @Router /admin/orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {

	//get pagination params
	page := c.Query("p")
	limit := c.Query("l")
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	// Validate and convert the string parameters to integers
	pageInt, err1 := requestValidation.ValidateAndParseInt(page)
	limitInt, err2 := requestValidation.ValidateAndParseInt(limit)
	if err1 != nil || err2 != nil {
		fmt.Println("error parsing p parameter. error:", err1, ",", err2)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error in url. Try Again", fmt.Errorf("%v, %v", err1, err2)))
		return
	}

	//get orders
	var orders *[]response.ResponseOrderInfo
	orders, message, err := h.orderUseCase.GetOrders(pageInt, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.GetOrdersResponse{
		Status:     "success",
		Message:    "Orders fetched successfully",
		OrdersInfo: *orders,
	})
}

// cancel order of user
// @Summary Cancel Order
// @Description User can cancel an order which is not yet delivered
// @Tags User/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param cancelOrderReq body requestModels.CancelOrderReq true "Cancel Order Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /cancelorder [patch]
func (h *OrderHandler) CancelMyOrder(c *gin.Context) {

	//get req from body
	var req *requestModels.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error cancelling order. Try Again", err))
		return
	}

	//cancel order
	message, err := h.orderUseCase.CancelOrderByUser(req.OrderID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Order cancelled successfully",
	})

}

// cancel order of any user with userID by admin
// @Summary Cancel Order
// @Description Admin can cancel an order which is not yet delivered
// @Tags Admin/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param cancelOrderReq body requestModels.CancelOrderReq true "Cancel Order Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/cancelorder [patch]
func (h *OrderHandler) CancelOrderByAdmin(c *gin.Context) {

	//get req from body
	var req *requestModels.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//cancel order
	message, err := h.orderUseCase.CancelOrderByAdmin(req.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Order cancelled successfully",
	})
}

// return order of user
// @Summary Return Order
// @Description User can request for returning an order which is already delivered
// @Tags User/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param returnOrderReq body requestModels.ReturnOrderReq true "Return Order Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /returnorder [patch]
func (h *OrderHandler) ReturnMyOrder(c *gin.Context) {

	//get req from body
	var req *requestModels.ReturnOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error returning order. Try Again", err))
		return
	}

	//return order
	message, err := h.orderUseCase.ReturnOrderRequestByUser(req.OrderID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: message,
	})
}

// mark order as returned by admin
// @Summary Mark Order as Returned
// @Description Admin can mark an order as returned when it is returned by the user and received by the admin
// @Tags Admin/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param orderID query string true "Order ID"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/markorderasreturned [patch]
func (h *OrderHandler) MarkOrderAsReturned(c *gin.Context) {

	//get req from body
	var req *requestModels.ReturnOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//mark order as returned
	message, err := h.orderUseCase.MarkOrderAsReturned(req.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Order marked as returned successfully"))
}

// MarkOrderAsDelivered by admin
// @Summary Mark Order as Delivered
// @Description Admin can mark an order as delivered when it is delivered to the user
// @Tags Admin/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param orderID query string true "Order ID"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/markdelivery [patch]
func (h *OrderHandler) MarkOrderAsDelivered(c *gin.Context) {

	//get req from body
	var req *requestModels.MarkOrderAsDeliveredReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//mark order as delivered
	message, err := h.orderUseCase.MarkOrderAsDelivered(req.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Order marked as delivered successfully"))
}

// get invoice
// @Summary Get Invoice
// @Description Get Invoice of an order
// @Tags User/Order
// @Produce json
// @Security BearerTokenAuth
// @Param orderID query string true "Order ID"
// @Success 200 {file} application/pdf
// @Failure 400 {object} response.SME{}
// @Router /order-invoice [get]
func (h *OrderHandler) GetInvoiceOfOrder(c *gin.Context) {

	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting invoice. Try Again", err))
		return
	}

	orderIdParam := c.Query("orderID")
	orderId, err := strconv.Atoi(orderIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error getting invoice. Try Again", err))
		return
	}

	//get invoice
	invoiceURL, message, err := h.orderUseCase.GetInvoiceOfOrder(userID, uint(orderId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	// c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	// c.Header("Content-Type", "application/pdf") // Set the correct content type for PDF
	// c.Redirect(http.StatusFound, *invoiceURL)
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.File(*invoiceURL)
}

// GetAddressForCheckout
// @Summary Get Address for Checkout
// @Description Get Address for Checkout
// @Tags User/Cart
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetAddressesForCheckoutResponse
// @Failure 400 {object} response.SME{}
// @Router /selectaddress [get]
func (h *OrderHandler) GetAddressForCheckout(c *gin.Context) {

	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting address for checkout. Try Again", err))
		return
	}

	//get address for checkout
	address, totalQuantiy, totalValue, message, err := h.orderUseCase.GetAddressForCheckout(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.GetAddressesForCheckoutResponse{
		Status:       "success",
		Message:      "Address fetched successfully",
		Addresses:    *address,
		TotalQuantiy: totalQuantiy,
		TotalValue:   totalValue,
	})
}

// SetAddressGetCoupons
// @Summary Set Address and Get Coupons
// @Description Set Address and Get Coupons
// @Tags User/Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param setAddGetCouponsReq body requestModels.SetAddressForCheckOutReq true "Set Address and Get Coupons Request"
// @Success 200 {object} response.SetAddrGetCouponsResponse
// @Failure 400 {object} response.SME{}
// @Router /setaddr-selectcoupon [post]
func (h *OrderHandler) SetAddressGetCoupons(c *gin.Context) {

	//get req from body
	var req *requestModels.SetAddressForCheckOutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error setting address. Try Again", err))
		return
	}

	//set address and get coupons
	resp, message, err := h.orderUseCase.SetAddressGetCoupons(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SetCouponGetPaymentMethods
// @Summary Set Coupon and Get Payment Methods
// @Description Set Coupon and Get Payment Methods
// @Tags User/Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param setCouponGetPaymentMethodsReq body requestModels.SetCouponForCheckoutReq true "Set Coupon and Get Payment Methods Request"
// @Success 200 {object} response.GetPaymentMethodsForCheckoutResponse
// @Failure 400 {object} response.SME{}
// @Router /setcoupon-getpaymentmethods [post]
func (h *OrderHandler) SetCouponGetPaymentMethods(c *gin.Context) {

	//get req from body
	var req *requestModels.SetCouponForCheckoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error setting coupon. Try Again", err))
		return
	}

	//set coupon and get payment methods
	resp, message, err := h.orderUseCase.SetCouponGetPaymentMethods(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
