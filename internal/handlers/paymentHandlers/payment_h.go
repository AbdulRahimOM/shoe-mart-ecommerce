package paymentHandlers

import (
	"MyShoo/internal/config"
	e "MyShoo/internal/domain/customErrors"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecase "MyShoo/internal/usecase/interface"
	htmlRender "MyShoo/pkg/htmlTemplateRender"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase usecase.IPaymentUC
}

func NewPaymentHandler(paymentUseCase usecase.IPaymentUC) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// ProceedToPayment
// @Summary Proceed to payment
// @Description Proceed to payment
// @Tags User/Payment
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param proceedToPaymentReq body req.ProceedToPaymentReq{} true "Proceed to Payment Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /payment [post]
func (h *PaymentHandler) ProceedToPayViaRazorPay(c *gin.Context) {

	//get req from body
	var req request.ProceedToPaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	c.HTML(http.StatusOK, "payment.html", req)
	// c.JSON(http.StatusOK, req)

	//Rendering HTML for viewing payment (for testing)	(#dev mode)
	if config.IsLocalHostMode && config.ShouldRenderPaymentPage {
		renderPaymentPageInTestkit(req)
	}

}

// VerifyPayment
// @Summary Verify payment
// @Description Verify payment
// @Tags User/Payment
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param razorpay_payment_id formData string true "Razorpay Payment ID"
// @Param razorpay_order_id formData string true "Razorpay Order ID"
// @Param razorpay_signature formData string true "Razorpay Signature"
// @Success 200 {object} response.PaidOrderResponse{}
// @Failure 400 {object} response.SME{}
// @Router /payment/verify [post]
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {

	if err := c.Request.ParseForm(); err != nil {
		c.JSON(500, response.FailedSME("Error parsing form data", err))
		return
	}

	req := request.VerifyPaymentReq{
		RazorpayPaymentID: string(c.Request.Form.Get("razorpay_payment_id")),
		RazorpayOrderID:   string(c.Request.Form.Get("razorpay_order_id")),
		RazorpaySignature: string(c.Request.Form.Get("razorpay_signature")),
	}

	paymentValid, orderDetails, message, err := h.paymentUseCase.VerifyPayment(&req)
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
// @Summary Retry payment
// @Description Retry payment
// @Tags User/Payment
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.RetryPaymentReq{} true "Retry Payment Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /retrypayment [post]
func (h *PaymentHandler) RetryPayment(c *gin.Context) {

	//get req from body
	var req request.RetryPaymentReq
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
		c.JSON(http.StatusBadRequest, response.FailedSME("Error getting userID from token", err))
		return
	}

	//do the retry payment
	paymentInfo, message, err := h.paymentUseCase.RetryPayment(&req, userID)
	if err != nil {
		c.JSON(500, response.FailedSME(message, err))
		return
	}

	c.HTML(http.StatusOK, "payment.html", paymentInfo)

	//Rendering HTML for viewing payment (for testing)(#dev mode)
	if config.IsLocalHostMode && config.ShouldRenderPaymentPage {
		renderPaymentPageInTestkit(req)
	}
}
func renderPaymentPageInTestkit(req interface{}) { //(#dev mode)
	templatePath := filepath.Join(config.ExecutableDir, "internal/templates/payment.html")
	renderOutputPath := filepath.Join(config.ExecutableDir, "testKit/paymentOutput.html")
	err := htmlRender.RenderHTMLFromTemplate(templatePath, req, renderOutputPath)
	if err != nil {
		fmt.Println("Page loaded successfully. But, Coulnot produce testKit/paymentOutput.html file as rendered version. Go for alternative ways")
	} else {
		fmt.Println("Page loaded successfully. testKit/paymentOutput.html file produced as rendered version.")
	}
}
