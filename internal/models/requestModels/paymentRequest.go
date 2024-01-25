package requestModels

// make req as like ProceedToPaymentInfo
type ProceedToPaymentReq struct {
	PaymentKey         string  `json:"payment_key" validate:"required"`
	PaymentOrderID     string  `json:"payment_order_id" validate:"required"`
	OrderRefNo         string  `json:"order_ref_no" validate:"required"`
	TotalAmount        float32 `json:"total_amount" validate:"required,number"`
	Discount           float32 `json:"discount" validate:"number"`
	ShippingCharge     float32 `json:"shipping_charge" validate:"number"`
	TotalPayableAmount float32 `json:"total_payable_amount" validate:"required,number"`
	FirstName          string  `json:"first_name" validate:"required"`
	Email              string  `json:"email" validate:"required,email"`
	Phone              string  `json:"phone" validate:"required"`
}
