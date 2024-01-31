package paymentHandlers

import (
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	htmlRender "MyShoo/pkg/htmlTemplateRender"
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
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		// c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(paymentReq); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
		// c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	c.HTML(http.StatusOK, "payment.html", paymentReq)
	// c.JSON(http.StatusOK, retryReq)

	//Rendering HTML for viewing payment (for testing)
	err := htmlRender.RenderHTMLFromTemplate("internal/view/payment.html", paymentReq, "testKit/paymentOutput.html")
	if err != nil {
		fmt.Println("Page loaded successfully. But, Coulnot produce testKit/paymentOutput.html file as rendered version. Go for alternative ways")
	} else {
		fmt.Println("Page loaded successfully. testKit/paymentOutput.html file produced as rendered version.")
	}

}

// VerifyPayment
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	fmt.Println("Handler ::: verify payment handler")

	if err := c.Request.ParseForm(); err != nil {
		fmt.Println("Error parsing form data:", err)
		c.JSON(500, response.FailedSME("Error parsing form data", err))
		return
	}

	var request requestModels.VerifyPaymentReq = requestModels.VerifyPaymentReq{
		RazorpayPaymentID: string(c.Request.Form.Get("razorpay_payment_id")),
		RazorpayOrderID:   string(c.Request.Form.Get("razorpay_order_id")),
		RazorpaySignature: string(c.Request.Form.Get("razorpay_signature")),
	}

	paymentValid, orderDetails, message, err := h.paymentUseCase.VerifyPayment(&request)
	if err != nil {
		fmt.Println("Error verifying payment:", err)
		c.JSON(500, response.FailedSME("Error verifying payment", err))
		return
	}
	if !paymentValid {
		c.JSON(http.StatusExpectationFailed, response.FailedSME("Payment failed", nil))
		return
	}

	c.JSON(http.StatusOK, response.PaidOrderResponse{
		Status:    "success",
		Message:   message,
		OrderInfo: *orderDetails,
	})
}

// retry payment
func (h *PaymentHandler) RetryPayment(c *gin.Context) {

	//get req from body
	var retryReq requestModels.RetryPaymentReq
	if err := c.ShouldBindJSON(&retryReq); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(retryReq); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error getting userID from token", err))
		return
	}

	//do the retry payment
	paymentReq, message, err := h.paymentUseCase.RetryPayment(&retryReq, userID)
	if err != nil {
		c.JSON(500, response.FailedSME(message, err))
		return
	}

	c.HTML(http.StatusOK, "payment.html", paymentReq)

	//Rendering HTML for viewing payment (for testing)
	err = htmlRender.RenderHTMLFromTemplate("internal/view/payment.html", paymentReq, "testKit/paymentOutput.html")
	if err != nil {
		fmt.Println("Page loaded successfully. But, Coulnot produce testKit/paymentOutput.html file as rendered version. Go for alternative ways")
	} else {
		fmt.Println("Page loaded successfully. testKit/paymentOutput.html file produced as rendered version.")
	}
}
