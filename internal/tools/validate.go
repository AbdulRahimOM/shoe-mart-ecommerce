package tools

import "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"

func IsValidPaymentMethod(paymentMethod string) bool {
	//validate payment method
	for _, v := range entities.PaymentMethod {
		if v == paymentMethod {
			return true
		}
	}
	return false
}
