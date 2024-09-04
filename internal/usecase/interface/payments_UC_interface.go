package usecase

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
)

type IPaymentUC interface {
	VerifyPayment(req *request.VerifyPaymentReq) (bool, *entities.Order, *e.Error)
	RetryPayment(req *request.RetryPaymentReq, userID uint) (*request.ProceedToPaymentReq, *e.Error)
}
