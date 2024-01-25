package paymentHandlers

import (
	myshoo "MyShoo"
	"MyShoo/internal/models/requestModels"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	// paymentUseCase usecaseInterface.IPaymentUC
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		// paymentUseCase: paymentUseCase,
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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error, if any
		fmt.Println("Error reading request body:", err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Print the request body to the console
	fmt.Println("Request Body:", string(body))
	//get req from body
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.HTML(http.StatusBadRequest, "test.html", gin.H{
			"message": "Error happened with request. Try Again",
			"error":   err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "test.html", request)
}
