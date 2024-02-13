package usecase

import (
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
)

type IPaymentUC interface {
	VerifyPayment(req *request.VerifyPaymentReq) (bool, *entities.Order, string, error)
	RetryPayment(req *request.RetryPaymentReq, userID uint) (*request.ProceedToPaymentReq, string, error)
}
