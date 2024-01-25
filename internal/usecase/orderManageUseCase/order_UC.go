package orderManageUseCase

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repository_interface "MyShoo/internal/repository/interface"
	"MyShoo/internal/services"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

type OrderUseCase struct {
	userRepo    repository_interface.IUserRepo
	orderRepo   repository_interface.IOrderRepo
	cartRepo    repository_interface.ICartRepo
	productRepo repository_interface.IProductsRepo
}

func NewOrderUseCase(
	userRepo repository_interface.IUserRepo,
	orderRepo repository_interface.IOrderRepo,
	cartRepo repository_interface.ICartRepo,
	productRepo repository_interface.IProductsRepo,
) usecaseInterface.IOrderUC {
	return &OrderUseCase{
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *OrderUseCase) MakeOrder(req *requestModels.MakeOrderReq) (*entities.OrderInfo, *response.ProceedToPaymentInfo, string, error) {
	var orderInfo entities.OrderInfo
	var order entities.Order
	var orderItems []entities.OrderItem
	var totalCartOriginalAmount float32
	var message string
	var responseOrder *entities.Order

	//validate payment method
	paymentValid := tools.IsValidPaymentMethod(req.PaymentMethod)
	if !paymentValid {
		fmt.Println("Invalid payment method")
		return &orderInfo, nil, "Invalid payment method", errors.New("invalid payment method")
	}
	//validate coupon //need update

	//check if address exists
	var addressExists bool
	var err error
	addressExists, err = uc.userRepo.DoAddressExistsByIDForUser(req.AddressID, req.UserID)
	if err != nil {
		fmt.Println("Error occured while checking if address exists")
		return &orderInfo, nil, "Some error occured.", err
	}

	if !addressExists {
		return &orderInfo, nil, "Invalid address ID.", errors.New("address doesn't exist by ID")
	}
	//check if cart is empty
	cartEmpty, err := uc.cartRepo.IsCartEmpty(req.UserID)
	if err != nil {
		fmt.Println("Error occured while checking if cart is empty")
		return &orderInfo, nil, "Some error occured.", err
	}
	if cartEmpty {
		return &orderInfo, nil, "Cart is empty", errors.New("cart is empty")
	}
	//get cart
	var cart *[]entities.Cart
	cart, err = uc.cartRepo.GetCart(req.UserID)
	if err != nil {
		fmt.Println("Error occured while getting cart")
		return &orderInfo, nil, "Some error occured.", err
	}

	for _, cartItem := range *cart {
		//check for stock availability
		var stock uint
		stock, err = uc.productRepo.GetStockOfProduct(cartItem.ProductID)
		if err != nil {
			fmt.Println("Error occured while getting stock")
			return &orderInfo, nil, "Error occured while getting stock", err
		}
		if stock < cartItem.Quantity {
			message := fmt.Sprint("Stock not available for product with product ID:", cartItem.ProductID, ". Available stock left:", stock)
			return &orderInfo, nil, message, errors.New("stock not available")
		}

		//append cart's productID and Quantity to orderItems
		var orderItem entities.OrderItem
		orderItem.ProductID = cartItem.ProductID
		orderItem.Quantity = cartItem.Quantity
		orderItem.SalePriceOnOrder, err = uc.productRepo.GetPriceOfProduct(cartItem.ProductID)
		if err != nil {
			return &orderInfo, nil, "Error occured while getting price", err
		}
		orderItems = append(orderItems, orderItem)
		totalCartOriginalAmount += orderItem.SalePriceOnOrder * float32(orderItem.Quantity)

		var pq entities.PQ //PQ indicates product-quantity pairs
		pq.ProductID = cartItem.ProductID
		pq.Quantity = cartItem.Quantity
		orderInfo.OrderItems = append(orderInfo.OrderItems, pq)

	}

	//define/get order fields
	order.ReferenceNo, err = tools.MakeRandomUUID()
	if err != nil {
		fmt.Println("Error occured while generating reference number")
		return &orderInfo, nil, "Some error occured.", err
	}
	order.OrderDateAndTime = time.Now()
	order.UserID = req.UserID
	order.OriginalAmount = totalCartOriginalAmount
	order.CouponDiscount = 0                 //discount not yet ready. //need update
	order.FinalAmount = order.OriginalAmount //discount not yet ready. //need update
	// order.CouponID = req.CouponID		//not yet ready. //need update
	order.PaymentMethod = req.PaymentMethod
	order.PaymentStatus = "not paid"
	order.AddressID = req.AddressID
	if strings.ToUpper(req.PaymentMethod) == "COD" {
		order.Status = "placed"
	} else {
		order.Status = "payment pending"
		order.TransactionID, err = services.CreateRazorpayOrder(order)
		if err != nil {
			return &orderInfo, nil, "Some error occured.", err
		}
	}
	if order.Status == "payment pending" {
		message = "Order placed successfully. "
		responseOrder, err = uc.orderRepo.MakeOrder_UpdateStock_ClearCart(&order, &orderItems)
	} else {
		message = "Proceed to payment"
		responseOrder, err = uc.orderRepo.MakeOrder(&order, &orderItems)
	}
	//make order
	
	
	if err != nil {
		fmt.Println("Error occured while placing order")
		return &orderInfo, nil, "Error occured while placing order. Try again or", err
	}

	orderInfo.OrderDetails = *responseOrder

	proceedToPaymentInfo := response.ProceedToPaymentInfo{
		PaymentKey:         os.Getenv("RAZORPAY_KEY_ID"),
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


	return &orderInfo, &proceedToPaymentInfo, message, nil
}

// GetOrdersOfUser
func (uc *OrderUseCase) GetOrdersOfUser(userID uint, page int, limit int) (*[]response.ResponseOrderInfo, string, error) {
	var orders *[]entities.DetailedOrderInfo
	var responseOrders []response.ResponseOrderInfo
	var err error
	offset := (page - 1) * limit
	orders, err = uc.orderRepo.GetOrdersOfUser(userID, offset, limit)
	if err != nil {
		fmt.Println("Error occured while getting orders")
		return &responseOrders, "Error occured while getting orders", err
	}

	//convert orders to responseOrders using "github.com/jinzhu/copier"
	err = copier.Copy(&responseOrders, &orders)
	if err != nil {
		fmt.Println("Error occured while copying orders to responseOrders")
		return &responseOrders, "Error occured while copying orders to responseOrders", err
	}

	return &responseOrders, "Orders fetched successfully", nil
}

// GetAllOrders
func (uc *OrderUseCase) GetOrders(page int, limit int) (*[]response.ResponseOrderInfo, string, error) {
	var orders *[]entities.DetailedOrderInfo
	var responseOrders []response.ResponseOrderInfo
	var err error
	offset := (page - 1) * limit
	orders, err = uc.orderRepo.GetOrders(offset, limit)
	if err != nil {
		fmt.Println("Error occured while getting orders")
		return &responseOrders, "Error occured while getting orders", err
	}

	//convert orders to responseOrders using "github.com/jinzhu/copier"
	err = copier.Copy(&responseOrders, &orders)
	if err != nil {
		fmt.Println("Error occured while copying orders to responseOrders")
		return &responseOrders, "Error occured while copying orders to responseOrders", err
	}
	return &responseOrders, "Orders fetched successfully", nil
}

// CancelOrder(orderID uint) (string, error)
func (uc *OrderUseCase) CancelOrderByUser(orderID uint, userID uint) (string, error) {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return "Some error occured.", err
	}
	if !orderExists {
		return "Invalid/Corrupt request", errors.New("order doesn't exist by ID")
	}

	//check if order belongs to userID
	userIDFromOrder, err := uc.orderRepo.GetUserIDByOrderID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting userID")
		return "Some error occured.", err
	}

	//check if userID in argument and userID from order match
	if userID != userIDFromOrder {
		return "Corrupt request", errors.New("order doesn't belong to user")
	}

	//check if order is already cancelled
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return "Some error occured.", err
	}

	if orderStatus != "placed" {
		if orderStatus == "cancelled" {
			message := "Order is already in '" + orderStatus + "' status"
			return message, errors.New(message)
		} else {
			message := "Cannot cancel. Order is in '" + orderStatus + "' status"
			return message, errors.New(message)
		}
	}

	//cancel order, update stock
	err = uc.orderRepo.CancelOrder(orderID)
	if err != nil {
		fmt.Println("Error occured while cancelling order")
		return "Some error occured.", err
	}

	return "Order cancelled successfully", nil
}

// CancelOrderByAdmin(orderID uint) (string, error)
func (uc *OrderUseCase) CancelOrderByAdmin(orderID uint) (string, error) {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return "Some error occured.", err
	}
	if !orderExists {
		return "Corrupt request. Corrupt request. Invalid order ID.", errors.New("order doesn't exist by ID")
	}

	//check if order is already cancelled
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return "Some error occured.", err
	}
	if orderStatus != "placed" {
		message := "Cannot return. Order is in '" + orderStatus + "' status"
		fmt.Println(message)
		return message, errors.New(message)
	}

	//cancel order
	err = uc.orderRepo.CancelOrder(orderID)
	if err != nil {
		fmt.Println("Error occured while cancelling order")
		return "Some error occured.", err
	}

	return "Order cancelled successfully", nil
}

// ReturnOrderRequestByUser
func (uc *OrderUseCase) ReturnOrderRequestByUser(orderID, userID uint) (string, error) {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return "Some error occured.", err
	}
	if !orderExists {
		return "Corrupt request/Corrupt request. Invalid order ID.", errors.New("order doesn't exist by ID")
	}

	//check if order belongs to userID
	userIDFromOrder, err := uc.orderRepo.GetUserIDByOrderID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting userID")
		return "Some error occured.", err
	}
	if userID != userIDFromOrder {
		return "Order doesn't belong to user", errors.New("order doesn't belong to user")
	}

	//check if order is in "delivered" status
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return "Some error occured.", err
	}
	if orderStatus != "delivered" {
		message := "Cannot return. Order is in '" + orderStatus + "' status"
		fmt.Println(message)
		return message, errors.New(message)
	}

	//return order
	err = uc.orderRepo.ReturnOrderRequest(orderID)
	if err != nil {
		fmt.Println("Error occured while returning order")
		return "Some error occured.", err
	}

	return "Request for return submitted successfully", nil
}

// MarkOrderAsReturned
func (uc *OrderUseCase) MarkOrderAsReturned(orderID uint) (string, error) {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return "Some error occured.", err
	}
	if !orderExists {
		return "Corrupt request/Corrupt request. Invalid order ID.", errors.New("order doesn't exist by ID")
	}

	//check if order is in "return requested" status
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return "Some error occured.", err
	}
	if orderStatus != "return requested" {
		message := "Cannot mark as returned. Order is in '" + orderStatus + "' status"
		fmt.Println(message)
		return message, errors.New(message)
	}

	//mark order as returned
	err = uc.orderRepo.MarkOrderAsReturned(orderID)
	if err != nil {
		fmt.Println("Error occured while marking order as returned")
		return "Some error occured.", err
	}

	return "Order marked as returned successfully", nil
}

// MarkOrderAsDelivered
func (uc *OrderUseCase) MarkOrderAsDelivered(orderID uint) (string, error) {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return "Some error occured.", err
	}
	if !orderExists {
		return "Corrupt request. Invalid order ID.", errors.New("order doesn't exist by ID")
	}

	//check if order is in "placed" status
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return "Some error occured.", err
	}
	if orderStatus != "placed" {
		if orderStatus == "delivered" {
			message := "Order is already in '" + orderStatus + "' status"
			return message, errors.New(message)
		} else {
			message := "Cannot mark as delivered. Order is in '" + orderStatus + "' status"
			return message, errors.New(message)
		}
	}

	//mark order as delivered
	err = uc.orderRepo.MarkOrderAsDelivered(orderID)
	if err != nil {
		fmt.Println("Error occured while marking order as delivered")
		return "Some error occured.", err
	}

	return "Order marked as delivered successfully", nil
}
