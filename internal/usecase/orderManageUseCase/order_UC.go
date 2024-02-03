package orderManageUseCase

import (
	"MyShoo/internal/domain/config"
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repository_interface "MyShoo/internal/repository/interface"
	"MyShoo/internal/services"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	myMath "MyShoo/pkg/math"
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
	var couponDiscount, totalProductValue float32
	var orderItems []entities.OrderItem
	var status, message string
	var transactionID string

	//validate payment method
	paymentValid := tools.IsValidPaymentMethod(req.PaymentMethod)
	if !paymentValid {
		return nil, nil, "Invalid payment method", errors.New("invalid payment method")
	}

	//check if address exists
	addressExists, err := uc.userRepo.DoAddressExistsByIDForUser(req.AddressID, req.UserID)
	if err != nil {
		fmt.Println("Error occured while checking if address exists")
		return nil, nil, "Some error occured.", err
	}

	if !addressExists {
		return nil, nil, "Invalid address ID.", errors.New("address doesn't exist by ID")
	}
	address, err := uc.userRepo.GetUserAddress(req.AddressID)
	if err != nil {
		return nil, nil, "Some error occured.", err
	}

	//check if cart is empty
	cartEmpty, err := uc.cartRepo.IsCartEmpty(req.UserID)
	if err != nil {
		fmt.Println("Error occured while checking if cart is empty")
		return nil, nil, "Some error occured.", err
	}
	if cartEmpty {
		return nil, nil, "Cart is empty", errors.New("cart is empty")
	}
	//get cart
	var cart *[]entities.Cart
	cart, totalProductValue, err = uc.cartRepo.GetCart(req.UserID)
	if err != nil {
		fmt.Println("Error occured while getting cart")
		return nil, nil, "Some error occured.", err
	}

	for _, cartItem := range *cart {
		//check for stock availability
		stock, err := uc.productRepo.GetStockOfProduct(cartItem.ProductID)
		if err != nil {
			fmt.Println("Error occured while getting stock")
			return nil, nil, "Error occured while getting stock", err
		}
		if stock < cartItem.Quantity {
			message := fmt.Sprint("Stock not available for product with product ID:", cartItem.ProductID, ". Available stock left:", stock)
			return nil, nil, message, errors.New("stock not available") //update needed
		}

		orderInfo.OrderItems = append(orderInfo.OrderItems, entities.PQ{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		})
	}

	referenceNo, err := tools.MakeRandomUUID()
	if err != nil {
		fmt.Println("Error occured while generating reference number")
		return nil, nil, "Some error occured.", err
	}

	shippingCharges := getShippingCharge(address, totalProductValue)

	if req.CouponID != 0 {
		coupon, message, err := uc.orderRepo.GetCouponByID(req.CouponID)
		if err != nil {
			return nil, nil, message, err
		}

		couponUsageCount, message, err := uc.orderRepo.GetCouponUsageCount(req.UserID, req.CouponID)
		if err != nil {
			return nil, nil, message, err
		}

		couponDiscount, message, err = getCouponDiscount(coupon, totalProductValue, couponUsageCount)
		if err != nil {
			return nil, nil, message, err
		}
	}

	finalAmount := totalProductValue + shippingCharges - couponDiscount

	if strings.ToUpper(req.PaymentMethod) == "COD" {
		if !config.DeliveryConfig.CashOnDeliveryAvailable {
			return nil, nil, msg.Forbidden, e.ErrCODNotAvailable
		}
		if totalProductValue > config.DeliveryConfig.MaxOrderAmountForCOD {
			return nil, nil, msg.Forbidden, e.ErrOrderExceedsMaxAmountForCOD
		}
		status = "placed"
	} else {

		status = "payment pending"
		transactionID, err = services.CreateRazorpayOrder(finalAmount, referenceNo)
		if err != nil {
			return nil, nil, "Some error occured.", err
		}
	}
	order := entities.Order{
		ReferenceNo:      referenceNo,
		OrderDateAndTime: time.Now(),
		UserID:           req.UserID,
		// DeliveredDate:    "",
		OriginalAmount: totalProductValue,
		CouponDiscount: couponDiscount,
		ShippingCharge: shippingCharges,
		FinalAmount:    finalAmount,
		CouponID:       req.CouponID,
		PaymentMethod:  req.PaymentMethod,
		Status:         status,
		AddressID:      req.AddressID,
		PaymentStatus:  "not paid",
		TransactionID:  transactionID,
		FkAddress:      *address,
	}

	fmt.Println("order.finalAmount=", order.FinalAmount)

	if order.Status == "payment pending" {
		message = "Order placed successfully. "
		order.ID, err = uc.orderRepo.MakeOrder_UpdateStock_ClearCart(&order, &orderItems)
		if err != nil {
			fmt.Println("Error occured while placing order")
			return nil, nil, "Error occured while placing order. Try again or", err
		}
	} else {
		message = "Proceed to payment"
		order.ID, err = uc.orderRepo.MakeOrder(&order, &orderItems)
		if err != nil {
			fmt.Println("Error occured while placing order")
			return nil, nil, "Error occured while placing order. Try again or", err
		}
	}
	//make order

	orderInfo.OrderDetails = order

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

	if orderStatus != "placed" && orderStatus != "payment pending" {
		if orderStatus == "cancelled" {
			message := "Order is already in '" + orderStatus + "' status"
			return message, errors.New(message)
		} else {
			message := "Cannot cancel. Order is in '" + orderStatus + "' status"
			return message, errors.New(message)
		}
	}

	//cancel order, update stock, refund to wallet if already paid
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
	if orderStatus != "placed" && orderStatus != "payment pending" {
		if orderStatus == "cancelled" {
			message := "Order is already in '" + orderStatus + "' status"
			return message, errors.New(message)
		} else {
			message := "Cannot cancel. Order is in '" + orderStatus + "' status"
			return message, errors.New(message)
		}
	}

	//cancel order, update stock, refund to wallet if already paid
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

	//mark order as returned, update stock, refund to wallet
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

	err = uc.orderRepo.MarkOrderAsDelivered(orderID)
	if err != nil {
		fmt.Println("Error occured while marking order as delivered")
		return "Some error occured.", err
	}

	return "Order marked as delivered successfully", nil
}

// // GetCheckOutInfo
// func (uc *OrderUseCase) GetCheckOutInfo(userID uint) (*response.CheckOutInfo, string, error) {
// 	var checkOutInfo response.CheckOutInfo
// 	var err error
// 	var message string
// 	//get quantity and price of cart
// 	checkOutInfo.ItemCount, checkOutInfo.TotalValue, message, err = uc.cartRepo.GetQuantityAndPriceOfCart(userID)
// 	if err != nil {
// 		return nil, message, err
// 	}

// 	//get address
// 	var address *[]entities.UserAddress
// 	address, err = uc.userRepo.GetUserAddresses(userID)
// 	if err != nil {
// 		return nil, "Some error occured.", err
// 	}

// 	//get coupons
// 	coupons, message, err := uc.orderRepo.GetActiveCoupons()
// 	if err != nil {
// 		return nil, message, err
// 	}

// 	checkOutInfo.Addresses = *address
// 	checkOutInfo.Coupons = *coupons

// 	return &checkOutInfo, "Checkout info fetched successfully", nil

// }

// func (uc *OrderUseCase) GetCheckOutEstimate(userID uint, req *requestModels.GetCheckoutEstimateReq) (*response.CheckoutEstimateResponse, *[]string, string, error) {
// 	var estimateResponse response.CheckoutEstimateResponse
// 	var message string
// 	var usageCount uint
// 	var address *entities.UserAddress

// 	//check if address exists for user
// 	addressExists, err := uc.userRepo.DoAddressExistsByIDForUser(req.AddressID, userID)
// 	if err != nil {
// 		fmt.Println("Error occured while checking if address exists")
// 		return nil, nil, "Some error occured.", err
// 	} else if !addressExists {
// 		return nil, nil, msg.InvalidRequest, errors.New("address doesn't exist by ID")
// 	}

// 	//get coupon by id
// 	coupon, message, err := uc.orderRepo.GetCouponByID(req.CouponID)
// 	if err != nil {
// 		return nil, nil, message, err
// 	}

// 	if usageCount, message, err = uc.orderRepo.GetCouponUsageCount(userID, req.CouponID); err != nil {
// 		return nil, nil, message, err
// 	}

// 	fmt.Println("userID=", userID)

// 	//get quantity and price of cart
// 	_, estimateResponse.ProductsValue, message, err = uc.cartRepo.GetQuantityAndPriceOfCart(userID)
// 	if err != nil {
// 		return nil, nil, message, err
// 	}

// 	estimateResponse.Discount, message, err = getCouponDiscount(coupon, estimateResponse.ProductsValue, usageCount)
// 	if err != nil {
// 		return nil, nil, message, err
// 	}

// 	address, err = uc.userRepo.GetUserAddress(req.AddressID)
// 	if err != nil {
// 		return nil, nil, "Some error occured.", err
// 	}

// 	estimateResponse.ShippingCharge, message, err = getShippingCharge(address, estimateResponse.ProductsValue)
// 	if err != nil {
// 		return nil, nil, message, err
// 	}

// 	estimateResponse.GrandTotal = estimateResponse.ProductsValue - estimateResponse.Discount + estimateResponse.ShippingCharge

// 	estimateResponse.PaymentMethods = *getPaymentMethods() //check if getting all payment methods here

// 	return &estimateResponse, getPaymentMethods(), "Estimate fetched successfully", nil
// }

func getPaymentMethods(orderValue float32) (*[]string, bool, string) {
	var paymentMethods []string
	var codAvailability bool = true
	var codAvailabilityNote string //if COD not available
	for _, method := range entities.PaymentMethod {
		if method == "COD" {
			if !config.DeliveryConfig.CashOnDeliveryAvailable {
				codAvailability = false
				codAvailabilityNote = "COD not available"
				continue
			} else if orderValue > config.DeliveryConfig.MaxOrderAmountForCOD {
				codAvailability = false
				codAvailabilityNote = "COD not available for order amount greater than " + fmt.Sprint(config.DeliveryConfig.MaxOrderAmountForCOD)
				continue
			} else {
				paymentMethods = append(paymentMethods, method)
				codAvailability = true
			}
		}
	}

	return &paymentMethods, codAvailability, codAvailabilityNote
}

func getShippingCharge(address *entities.UserAddress, productsValue float32) float32 {
	var pincode uint = address.Pincode

	if productsValue >= config.DeliveryConfig.OrderAmountForFreeDelivery {
		return 0
	}

	for _, v := range config.DeliveryConfig.FreeDeliveryPincodeRanges {
		if pincode >= v.Start && pincode <= v.End {
			return 0
		}
	}

	for _, v := range config.DeliveryConfig.IntermediatePincodeRanges {
		if pincode >= v.Start && pincode <= v.End {
			return config.DeliveryConfig.IntermediateDeliveryCharge
		}
	}

	return config.DeliveryConfig.DistantDeliveryCharge
}

func getCouponDiscount(coupon *entities.Coupon, orderValue float32, usageCount uint) (float32, string, error) {
	var discount float32
	switch {
	case coupon.Blocked:
		return 0, msg.Forbidden, errors.New("coupon is blocked")
	case coupon.StartDate.After(time.Now()):
		return 0, msg.Forbidden, errors.New("coupon is not active yet")
	case coupon.EndDate.Before(time.Now()):
		return 0, msg.Forbidden, errors.New("coupon is expired")
	case coupon.MinOrderValue > orderValue:
		return 0, msg.Forbidden, errors.New("order amount is less than required for coupon")
	case usageCount >= coupon.UsageLimit:
		return 0, msg.Forbidden, errors.New("coupon usage limit exceeded")
	default:
		if coupon.Type == "fixed" {
			discount = myMath.RoundFloat32(coupon.MaxDiscount, 2)
		} else { //if coupon.Type == "percentage" {
			discount = max(orderValue*coupon.Percentage/100, coupon.MaxDiscount)
			discount = myMath.RoundFloat32(discount, 2)
		}
		return discount, "", nil
	}
}

// GetAddressForCheckout implements usecaseInterface.IOrderUC.
func (uc *OrderUseCase) GetAddressForCheckout(userID uint) (*[]entities.UserAddress, uint, float32, string, error) {
	quantity, totalValue, message, err := uc.cartRepo.GetQuantityAndPriceOfCart(userID)
	if err != nil {
		return nil, 0, 0, message, err
	}

	addresses, err := uc.userRepo.GetUserAddresses(userID)
	if err != nil {
		return nil, 0, 0, "Some error occured.", err
	}

	return addresses, quantity, totalValue, "Addresses fetched successfully", nil
}

// SetAddressGetCoupons implements usecaseInterface.IOrderUC.
func (uc *OrderUseCase) SetAddressGetCoupons(userID uint, req *requestModels.SetAddressForCheckOutReq) (*response.SetAddrGetCouponsResponse, string, error) {
	// var resp response.SetAddrGetCouponsResponse
	var message string
	var err error
	address, err := uc.userRepo.GetUserAddress(req.AddressID)
	if err != nil {
		return nil, "Some error occured.", err
	} else if address.UserID != userID {
		return nil, "Address doesn't belong to user", errors.New("address doesn't belong to user")
	}

	totalQuantiy, totalProductsValue, message, err := uc.cartRepo.GetQuantityAndPriceOfCart(userID)
	if err != nil {
		return nil, message, err
	}

	//get shipping charge
	shippingCharge := getShippingCharge(address, totalProductsValue)
	if err != nil {
		return nil, message, err
	}

	//get coupons
	coupons, message, err := uc.orderRepo.GetActiveCoupons()
	if err != nil {
		return nil, message, err
	}
	var respCoupons []response.ResponseCoupon
	err = copier.Copy(&respCoupons, &coupons)
	if err != nil {
		fmt.Println("Error occured while copying coupons to responseCoupons, error:", err)
		return nil, "Some error occured.", err
	}

	response := response.SetAddrGetCouponsResponse{
		Status:       "success",
		Message:      "Address and coupons fetched successfully",
		Coupons:      respCoupons,
		Address:      *address,
		TotalQuantiy: totalQuantiy,
		BillSumary: response.BillBeforeCoupon{
			TotalProductsValue: totalProductsValue,
			ShippingCharge:     shippingCharge,
		},
	}

	return &response, "", nil
}

// SetCouponGetPaymentMethods implements usecaseInterface.IOrderUC.
func (uc *OrderUseCase) SetCouponGetPaymentMethods(userID uint, req *requestModels.SetCouponForCheckoutReq) (*response.GetPaymentMethodsForCheckoutResponse, string, error) {

	address, err := uc.userRepo.GetUserAddress(req.AddressID)
	if err != nil {
		return nil, "Some error occured.", err
	} else if address.UserID != userID {
		return nil, "Address doesn't belong to user", errors.New("address doesn't belong to user")
	}

	//get coupon by id

	totalQuantiy, totalProductsValue, message, err := uc.cartRepo.GetQuantityAndPriceOfCart(userID)
	if err != nil {
		return nil, message, err
	}

	shippingCharge := getShippingCharge(address, totalProductsValue)
	if err != nil {
		return nil, message, err
	}

	coupon, message, err := uc.orderRepo.GetCouponByID(req.CouponID)
	if err != nil {
		return nil, message, err
	}

	usageCount, message, err := uc.orderRepo.GetCouponUsageCount(userID, req.CouponID)
	if err != nil {
		return nil, message, err
	}

	couponDiscount, message, err := getCouponDiscount(coupon, totalProductsValue, usageCount)
	if err != nil {
		return nil, message, err
	}

	walletBalance, err := uc.userRepo.GetWalletBalance(userID)
	if err != nil {
		return nil, "Some error occured.", err
	}

	grandTotal := totalProductsValue - couponDiscount + shippingCharge
	paymentMethods, codAvailability, codAvailabilityNote := getPaymentMethods(grandTotal)

	var respCoupon response.ResponseCoupon
	err = copier.Copy(&respCoupon, &coupon)
	if err != nil {
		fmt.Println("Error occured while copying coupon to responseCoupon")
		return nil, "Some error occured.", err
	}

	resp := response.GetPaymentMethodsForCheckoutResponse{
		Status:       "success",
		Message:      "Payment methods fetched successfully",
		Address:      *address,
		TotalQuantiy: totalQuantiy,
		BillSumary: response.BillAfterCoupon{
			TotalProductsValue: totalProductsValue,
			CouponApplied:      true,
			Coupon:             respCoupon,
			CouponDiscount:     couponDiscount,
			ShippingCharge:     shippingCharge,
			GrandTotal:         grandTotal,
		},
		CODAvailability:     codAvailability,
		CODAvailabilityNote: codAvailabilityNote,
		PaymentMethods:      *paymentMethods,
		WalletBalance:       walletBalance,
	}

	return &resp, "", nil
}

func (uc *OrderUseCase) GetInvoiceOfOrder(userID uint, orderID uint) (*string, string, error) {

	// //check if order exists
	// orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	// if err != nil {
	// 	fmt.Println("Error occured while checking if order exists")
	// 	return nil,"Some error occured.", err
	// }
	// if !orderExists {
	// 	return nil,"Corrupt request. Invalid order ID.", errors.New("order doesn't exist by ID")
	// }

	// //check if order belongs to userID
	// userIDFromOrder, err := uc.orderRepo.GetUserIDByOrderID(orderID)
	// if err != nil {
	// 	fmt.Println("Error occured while getting userID")
	// 	return nil,"Some error occured.", err
	// }

	// if userID != userIDFromOrder {
	// 	return nil,"Order doesn't belong to user", errors.New("order doesn't belong to user")
	// }

	// //check if payment status is "paid"
	// paymentStatus, err := uc.orderRepo.GetPaymentStatusByID(orderID)
	// if err != nil {
	// 	fmt.Println("Error occured while getting payment status")
	// 	return nil,"Some error occured.", err
	// }
	// if paymentStatus != "paid" {
	// 	message := "Cannot generate invoice. Payment status is '" + paymentStatus + "'"
	// 	fmt.Println(message)
	// 	return nil,message, errors.New(message)
	// }

	// //get orderInfo
	// var orderInfo *entities.OrderInfo
	// order, err := uc.orderRepo.GetOrderSummaryByID(orderID)
	// if err != nil {
	// 	fmt.Println("Error occured while getting order summary")
	// 	return nil,"Some error occured.", err
	// }

	panic("unimplemented")
}
