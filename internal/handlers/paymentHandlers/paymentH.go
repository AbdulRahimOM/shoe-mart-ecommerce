package paymentHandlers

import (
	myshoo "MyShoo"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase usecaseInterface.IPaymentUC
}

func NewPaymentHandler(paymentUseCase usecaseInterface.IPaymentUC) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// ProceedToPayment
func (h *PaymentHandler) ProceedToPayViaRazorPay(c *gin.Context) {
	fmt.Println("Handler ::: proceed to payment handler")

	//get req from body
	var paymentReq requestModels.ProceedToPaymentReq
	if err := c.ShouldBindJSON(&paymentReq); err != nil {
		c.HTML(http.StatusBadRequest, "payment.html", gin.H{
			"message": "Error happened with request. Try Again",
			"error":   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(paymentReq); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.HTML(http.StatusBadRequest, "payment.html", gin.H{
			"message": "Invalid request. Try again.",
			"error":   errResponse,
		})
		return
	}

	c.HTML(http.StatusOK, "payment.html", paymentReq)
}

// temporary GET method
func (h *PaymentHandler) ProceedToPayViaRazorPay2(c *gin.Context) {
	fmt.Println("Handler ::: proceed to payment handler")

	c.HTML(http.StatusOK, "payment.html", myshoo.ProceedToPaymentInfo)
}

// VerifyPayment
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	fmt.Println("Handler ::: verify payment handler")

	if err := c.Request.ParseForm(); err != nil {
		fmt.Println("Error parsing form data:", err)
		c.JSON(500, response.SME{
			Status:  "failed",
			Message: "Error parsing form data. Try Again",
			Error:   err.Error(),
		})
		return
	}

	var request requestModels.VerifyPaymentReq= requestModels.VerifyPaymentReq{
		RazorpayPaymentID: string(c.Request.Form.Get("razorpay_payment_id")),	
		RazorpayOrderID: string(c.Request.Form.Get("razorpay_order_id")),
		RazorpaySignature: string(c.Request.Form.Get("razorpay_signature")),
	}

	paymentValid,orderDetails, message, err := h.paymentUseCase.VerifyPayment(&request)
	if err != nil {
		fmt.Println("Error verifying payment:", err)
		c.JSON(500, response.SME{
			Status:  "failed",
			Message: message,
			Error:   err.Error(),
		})
		return
	}
	if !paymentValid {
		c.JSON(http.StatusExpectationFailed, response.SME{
			Status:  "failed",	
			Message: message,
			Error:   "Payment failed",
		})
		return
	}

	c.JSON(http.StatusOK, response.PaidOrderResponse{
		Status:  "success",
		Message: message,
		OrderInfo: *orderDetails,
	})

}
