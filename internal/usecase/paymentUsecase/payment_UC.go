package paymentusecase

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	repository_interface "MyShoo/internal/repository/interface"
	"MyShoo/internal/services"
	usecaseInterface "MyShoo/internal/usecase/interface"
)

type PaymentUC struct {
	paymentRepo repository_interface.IPaymentRepo
	orderRepo   repository_interface.IOrderRepo
	cartRepo    repository_interface.ICartRepo
	productRepo repository_interface.IProductsRepo
}

func NewPaymentUseCase(
	paymentRepo repository_interface.IPaymentRepo,
	orderRepo repository_interface.IOrderRepo,
	cartRepo repository_interface.ICartRepo,
	productRepo repository_interface.IProductsRepo,
) usecaseInterface.IPaymentUC {
	return &PaymentUC{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// VerifyPayment implements usecaseInterface.IPaymentUC.
func (uc *PaymentUC) VerifyPayment(req *requestModels.VerifyPaymentReq) (bool, *entities.Order, string, error) {
	//verify payment
	isPaymentValid:=services.VerifyPayment(req.RazorpayOrderID,req.RazorpayPaymentID,req.RazorpaySignature)
	if !isPaymentValid{
		return false, nil, "Payment failed", nil
	}

	//get orderID from orders table with transactionID=razorpayOrderID
	orderID, err := uc.orderRepo.GetOrderByTransactionID(req.RazorpayOrderID)
	if err != nil {
		return false, nil, "Some error occured while fetching order details", err
	}

	//UpdateOrderToPaid_UpdateStock_ClearCart
	order, err := uc.orderRepo.UpdateOrderToPaid_UpdateStock_ClearCart(orderID)
	if err != nil {
		return false, nil, "", err
	}

	return true, order, "Hoorray!!.. Payment recieved. Your order is place.", nil





}