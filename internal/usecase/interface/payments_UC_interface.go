package usecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
)

type IPaymentUC interface {
	VerifyPayment(req *request.VerifyPaymentReq) (bool, *entities.Order, *e.Error)
	RetryPayment(req *request.RetryPaymentReq, userID uint) (*request.ProceedToPaymentReq,  *e.Error)
}
