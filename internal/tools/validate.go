package tools

import "MyShoo/internal/domain/entities"

func IsValidPaymentMethod(paymentMethod string) bool {
	//validate payment method
	for _, v := range entities.PaymentMethod {
		if v == paymentMethod {
			return true
		}
	}
	return false
}
