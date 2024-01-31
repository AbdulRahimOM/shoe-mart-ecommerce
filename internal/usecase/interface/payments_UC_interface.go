package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
)

type IPaymentUC interface {
	VerifyPayment(req *requestModels.VerifyPaymentReq) (bool,*entities.Order, string, error)
	RetryPayment(req *requestModels.RetryPaymentReq, userID uint) (*requestModels.ProceedToPaymentReq, string, error)
}
