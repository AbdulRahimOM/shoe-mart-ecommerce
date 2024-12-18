package paymentusecase

import (
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/config"
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	repo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/services"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
)

type PaymentUC struct {
	orderRepo repo.IOrderRepo
}

func NewPaymentUseCase(
	orderRepo repo.IOrderRepo,
) usecase.IPaymentUC {
	return &PaymentUC{
		orderRepo: orderRepo,
	}
}

// VerifyPayment implements usecase.IPaymentUC.
func (uc *PaymentUC) VerifyPayment(req *request.VerifyPaymentReq) (bool, *entities.Order, *e.Error) {
	//verify payment
	isPaymentValid := services.VerifyPayment(req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	if !isPaymentValid {
		return false, nil, nil
	}

	//get orderID from orders table with transactionID=razorpayOrderID
	orderID, err := uc.orderRepo.GetOrderByTransactionID(req.RazorpayOrderID)
	if err != nil {
		return false, nil, err
	}

	//UpdateOrderToPaid_UpdateStock_ClearCart
	order, err := uc.orderRepo.UpdateOrderToPaid_UpdateStock_ClearCart(orderID)
	if err != nil {
		return false, nil, err
	}

	return true, order, nil
}

// RetryPayment(req *request.RetryPaymentReq) (*req.ProceedToPaymentReq, string, error)
func (uc *PaymentUC) RetryPayment(req *request.RetryPaymentReq, userID uint) (*request.ProceedToPaymentReq, *e.Error) {
	//check if order belongs to user
	userIDOfOrder, err := uc.orderRepo.GetUserIDByOrderID(req.OrderID)
	if err != nil {
		return nil, err
	} else {
		if userIDOfOrder != userID {
			return nil, e.SetError("OrderNotOfUser", nil, 400)
		}
	}

	//get order details
	order, err := uc.orderRepo.GetOrderSummaryByID(req.OrderID)
	if err != nil {
		return nil, err
	}

	//check if payment method is online and order is not paid
	if order.Status != "payment pending" {
		return nil, e.SetError("order payment status is not 'payment pending'", nil, 400)
	} else {
		var errr error
		if order.TransactionID, errr = services.CreateRazorpayOrder(order.FinalAmount, order.ReferenceNo); errr != nil {
			return nil, e.SetError("Service error at creating razorpay order:", errr, 500)
		}

		//update order with transactionID
		err = uc.orderRepo.UpdateOrderTransactionID(order.ID, order.TransactionID)
		if err != nil {
			return nil, err
		}
	}

	proceedToPaymentReq := request.ProceedToPaymentReq{
		PaymentKey:         config.RazorpayKeyId,
		PaymentOrderID:     order.TransactionID, //need update //payment-u
		OrderRefNo:         order.ReferenceNo,
		TotalAmount:        order.FinalAmount,
		Discount:           0,
		ShippingCharge:     0,
		TotalPayableAmount: order.FinalAmount,
		FirstName:          order.FkAddress.FirstName,
		Email:              order.FkAddress.Email,
		Phone:              order.FkAddress.Phone,
		CallBackURL:        config.BaseURL + "/payment/verify",
	}

	return &proceedToPaymentReq, nil
}
