package orderhandler

import (
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecase "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase usecase.IOrderUC
}

func NewOrderHandler(orderUseCase usecase.IOrderUC) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// MakeOrder
// @Summary Make Order
// @Description Make Order
// @Tags User/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param makeOrderReq body req.MakeOrderReq{} true "Make Order Request"
// @Success 201 {object} response.CODOrderResponse
// @Success 201 {object} response.OnlinePaymentOrderResponse
// @Failure 400 {object} response.SME{}
// @Router /makeorder [post]
func (h *OrderHandler) MakeOrder(c *gin.Context) {

	var req *request.MakeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}
	req.UserID = userID

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//make order
	orderInfo, proceedToPaymentInfo, err := h.orderUseCase.MakeOrder(req)
	if err != nil {
		switch err.Error() {
		case "cart is empty":
			c.JSON(http.StatusBadRequest, response.FromError(err))
		case "stock not available":
			c.JSON(http.StatusForbidden, response.FromError(err))
		default:
			c.JSON(http.StatusInternalServerError, response.FromError(err))
		}
		return

	}
	if proceedToPaymentInfo == nil {
		c.JSON(http.StatusCreated, response.CODOrderResponse{
			OrderInfo: *orderInfo,
		})
	} else {
		c.JSON(http.StatusCreated, response.OnlinePaymentOrderResponse{
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
	pageInt, errr := requestValidation.ValidateAndParseInt(page)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error parsing page(p) parameter. error:", errr))
		return
	}

	limitInt, errr := requestValidation.ValidateAndParseInt(limit)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error parsing limit(l) parameter. error:", errr))
		return
	}

	//get userID from token
	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//get orders
	orders, err := h.orderUseCase.GetOrdersOfUser(userID, pageInt, limitInt)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.GetOrdersResponse{
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
	pageInt, errr := requestValidation.ValidateAndParseInt(page)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error parsing page(p) parameter. error:", errr))
		return
	}

	limitInt, errr := requestValidation.ValidateAndParseInt(limit)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("error parsing limit(l) parameter. error:", errr))
		return
	}

	//get orders
	orders, err := h.orderUseCase.GetOrders(pageInt, limitInt)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.GetOrdersResponse{
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
// @Param cancelOrderReq body req.CancelOrderReq true "Cancel Order Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /cancelorder [patch]
func (h *OrderHandler) CancelMyOrder(c *gin.Context) {

	//get req from body
	var req *request.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//get userID from token
	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//cancel order
	err := h.orderUseCase.CancelOrderByUser(req.OrderID, userID)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, nil)

}

// cancel order of any user with userID by admin
// @Summary Cancel Order
// @Description Admin can cancel an order which is not yet delivered
// @Tags Admin/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param cancelOrderReq body req.CancelOrderReq true "Cancel Order Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/cancelorder [patch]
func (h *OrderHandler) CancelOrderByAdmin(c *gin.Context) {

	//get req from body
	var req *request.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//cancel order
	err := h.orderUseCase.CancelOrderByAdmin(req.OrderID)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// return order of user
// @Summary Return Order
// @Description User can request for returning an order which is already delivered
// @Tags User/Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param returnOrderReq body req.ReturnOrderReq true "Return Order Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /returnorder [patch]
func (h *OrderHandler) ReturnMyOrder(c *gin.Context) {

	//get req from body
	var req *request.ReturnOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//get userID from token
	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//return order
	err := h.orderUseCase.ReturnOrderRequestByUser(req.OrderID, userID)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
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
	var req *request.ReturnOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//mark order as returned
	err := h.orderUseCase.MarkOrderAsReturned(req.OrderID)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
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
	var req *request.MarkOrderAsDeliveredReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//mark order as delivered
	err := h.orderUseCase.MarkOrderAsDelivered(req.OrderID)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
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

	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	orderIdParam := c.Query("orderID")
	orderId, errr := strconv.Atoi(orderIdParam)
	if errr != nil {
		c.JSON(http.StatusBadRequest, response.MsgAndError("invalid order ID in query param", errr))
		return
	}

	//get invoice
	invoiceURL, err := h.orderUseCase.GetInvoiceOfOrder(userID, uint(orderId))
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
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

	userID, errr := tools.GetUserID(c)
	if errr != nil {
		// c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting address for checkout. Try Again", err))
		c.JSON(http.StatusInternalServerError, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//get address for checkout
	address, totalQuantiy, totalValue, err := h.orderUseCase.GetAddressForCheckout(userID)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.GetAddressesForCheckoutResponse{
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
// @Param setAddGetCouponsReq body req.SetAddressForCheckOutReq true "Set Address and Get Coupons Request"
// @Success 200 {object} response.SetAddrGetCouponsResponse
// @Failure 400 {object} response.SME{}
// @Router /setaddr-selectcoupon [post]
func (h *OrderHandler) SetAddressGetCoupons(c *gin.Context) {

	//get req from body
	var req *request.SetAddressForCheckOutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//get userID from token
	userID, errr := tools.GetUserID(c)
	if errr != nil {
		c.JSON(http.StatusForbidden, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//set address and get coupons
	resp, err := h.orderUseCase.SetAddressGetCoupons(userID, req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
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
// @Param setCouponGetPaymentMethodsReq body req.SetCouponForCheckoutReq true "Set Coupon and Get Payment Methods Request"
// @Success 200 {object} response.GetPaymentMethodsForCheckoutResponse
// @Failure 400 {object} response.SME{}
// @Router /setcoupon-getpaymentmethods [post]
func (h *OrderHandler) SetCouponGetPaymentMethods(c *gin.Context) {

	//get req from body
	var req *request.SetCouponForCheckoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//get userID from token
	userID, errr := tools.GetUserID(c)
	if errr != nil {
		// c.JSON(http.StatusInternalServerError, response.FailedSME("Error setting coupon. Try Again", err))
		c.JSON(http.StatusForbidden, response.MsgAndError("error getting user ID from token. error:", errr))
		return
	}

	//set coupon and get payment methods
	resp, err := h.orderUseCase.SetCouponGetPaymentMethods(userID, req)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
