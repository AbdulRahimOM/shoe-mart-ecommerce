package ordermanagementHandlers

import (
	myshoo "MyShoo"
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
func (h *OrderHandler) MakeOrder(c *gin.Context) {
	fmt.Println("Handler ::: make order handler")
	var req *requestModels.MakeOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error: %v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//check if userID in token and request body match
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error making order. Try Again", err))
		return
	}
	if userID != req.UserID {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error making order. Try Again",
			Error:   "User ID in token and request body do not match. Corrupted request!!",
		})
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
		myshoo.TemporarySolution(*proceedToPaymentInfo)
	}
}

// Get Orders of the user
func (h *OrderHandler) GetOrdersOfUser(c *gin.Context) {
	fmt.Println("Handler ::: get orders of user handler")

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
func (h *OrderHandler) GetOrders(c *gin.Context) {
	fmt.Println("Handler ::: get all orders handler")

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
func (h *OrderHandler) CancelMyOrder(c *gin.Context) {
	fmt.Println("Handler ::: cancel order handler")

	//get req from body
	var req *requestModels.CancelOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", fmt.Errorf(errResponse)))
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
func (h *OrderHandler) CancelOrderByAdmin(c *gin.Context) {

	//get req from body
	var req *requestModels.CancelOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", fmt.Errorf(errResponse)))
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
func (h *OrderHandler) ReturnMyOrder(c *gin.Context) {

	//get req from body
	var req *requestModels.ReturnOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", fmt.Errorf(errResponse)))
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

func (h *OrderHandler) MarkOrderAsReturned(c *gin.Context) {
	//get req from body
	var req *requestModels.ReturnOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", fmt.Errorf(errResponse)))
		return
	}

	//mark order as returned
	message, err := h.orderUseCase.MarkOrderAsReturned(req.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSME("Order marked as returned successfully"))
}

// MarkOrderAsDelivered
func (h *OrderHandler) MarkOrderAsDelivered(c *gin.Context) {

	//get req from body
	var req *requestModels.MarkOrderAsDeliveredReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", fmt.Errorf(errResponse)))
		return
	}

	//mark order as delivered
	message, err := h.orderUseCase.MarkOrderAsDelivered(req.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSME("Order marked as delivered successfully"))
}

// get invoice
func (h *OrderHandler) GetInvoiceOfOrder(c *gin.Context) {
	fmt.Println("Handler ::: 'GetInvoiceOfOrder' handler")

	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting invoice. Try Again", err))
		return
	}

	orderIdParam := c.Query("id")
	orderId, err := strconv.Atoi(orderIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error getting invoice. Try Again", err))
		return
	}

	//get invoice
	invoice, message, err := h.orderUseCase.GetInvoiceOfOrder(uint(orderId), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME(message, err))
		return
	}

	c.JSON(http.StatusOK, response.GetInvoiceResponse{
		Status:  "success",
		Message: "Invoice fetched successfully",
		Invoice: *invoice,
	})
}
