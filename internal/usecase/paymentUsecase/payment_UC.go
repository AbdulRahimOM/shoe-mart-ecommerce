package paymentusecase

import (
	"MyShoo/internal/config"
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	request "MyShoo/internal/models/requestModels"
	repo "MyShoo/internal/repository/interface"
	"MyShoo/internal/services"
	usecase "MyShoo/internal/usecase/interface"
	"errors"
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
			// return nil,  e.ErrOrderNotOfUser
			return nil, &e.Error{Err: errors.New("OrderNotOfUser"), StatusCode: 400}
		}
	}

	//get order details
	order, err := uc.orderRepo.GetOrderSummaryByID(req.OrderID)
	if err != nil {
		return nil, err
	}

	//check if payment method is online and order is not paid
	if order.Status != "payment pending" {
		return nil, &e.Error{Err: errors.New("order payment status is not 'payment pending'"), StatusCode: 400}
	} else {
		if order.TransactionID, err = services.CreateRazorpayOrder(order.FinalAmount, order.ReferenceNo);err != nil {
			return nil, &e.Error{Err: errors.New("Service error"), StatusCode: 500}
		}

		//update order with transactionID
		err = uc.orderRepo.UpdateOrderTransactionID(order.ID, order.TransactionID)
		if err != nil {
			return nil, "Error updating order with transactionID", err
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
	}

	return &proceedToPaymentReq, "Kinldy proceed to payment", nil
}
