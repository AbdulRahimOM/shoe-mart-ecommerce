package response

type ProceedToPaymentInfo struct {
	PaymentKey         string  `json:"payment_key" gorm:"column:payment_key"`
	PaymentOrderID     string  `json:"payment_order_id" gorm:"column:payment_order_id"`
	OrderRefNo         string  `json:"order_ref_no" gorm:"column:order_ref_no"`
	TotalAmount        float32 `json:"total_amount" gorm:"column:total_amount"`
	Discount           float32 `json:"discount" gorm:"column:discount"`
	ShippingCharge     float32 `json:"shipping_charge" gorm:"column:shipping_charge"`
	TotalPayableAmount float32 `json:"total_payable_amount" gorm:"column:total_payable_amount"`
	FirstName          string  `json:"first_name" gorm:"column:first_name"`
	Email              string  `json:"email" gorm:"column:email"`
	Phone              string  `json:"phone" gorm:"column:phone"`
}
