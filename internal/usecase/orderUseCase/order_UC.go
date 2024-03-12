package orderusecase

import (
	"MyShoo/internal/config"
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	repo "MyShoo/internal/repository/interface"
	"MyShoo/internal/services"
	"MyShoo/internal/tools"
	usecase "MyShoo/internal/usecase/interface"
	myMath "MyShoo/pkg/math"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf"
)

type OrderUseCase struct {
	userRepo    repo.IUserRepo
	orderRepo   repo.IOrderRepo
	cartRepo    repo.ICartRepo
	productRepo repo.IProductsRepo
}

func NewOrderUseCase(
	userRepo repo.IUserRepo,
	orderRepo repo.IOrderRepo,
	cartRepo repo.ICartRepo,
	productRepo repo.IProductsRepo,
) usecase.IOrderUC {
	return &OrderUseCase{
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *OrderUseCase) MakeOrder(req *request.MakeOrderReq) (*entities.OrderInfo, *response.ProceedToPaymentInfo, *e.Error) {
	var orderInfo entities.OrderInfo
	var couponDiscount, totalProductValue float32
	var orderItems []entities.OrderItem
	var transactionID string
	var status string

	//validate payment method
	paymentValid := tools.IsValidPaymentMethod(req.PaymentMethod)
	if !paymentValid {
		return nil, nil, &e.Error {Err:errors.New("invalid payment method"), StatusCode: 400}
	}

	//check if address exists
	addressExists, err := uc.userRepo.DoAddressExistsByIDForUser(req.AddressID, req.UserID)
	if err != nil {
		fmt.Println("Error occured while checking if address exists")
		return nil, nil, err
	}

	if !addressExists {
		return nil, nil, &e.Error {Err:errors.New("address doesn't exist by ID"), StatusCode: 400}
	}
	address, err := uc.userRepo.GetUserAddress(req.AddressID)
	if err != nil {
		return nil, nil, err
	}

	//check if cart is empty
	cartEmpty, err := uc.cartRepo.IsCartEmpty(req.UserID)
	if err != nil {
		fmt.Println("Error occured while checking if cart is empty")
		return nil, nil, err
	}
	if cartEmpty {
		return nil, nil, &e.Error {Err:errors.New("cart is empty"), StatusCode: 400}
	}
	//get cart
	var cart *[]entities.Cart
	cart, totalProductValue, err = uc.cartRepo.GetCart(req.UserID)
	if err != nil {
		return nil, nil, err
	}

	for _, cartItem := range *cart {
		//check for stock availability
		stock, err := uc.productRepo.GetStockOfProduct(cartItem.ProductID)
		if err != nil {
			fmt.Println("Error occured while getting stock")
			return nil, nil, err
		}
		if stock < cartItem.Quantity {
			err := fmt.Errorf("Stock not available for product with product ID:", cartItem.ProductID, ". Available stock left:", stock)
			return nil, nil, err //update needed
		}

		orderItems = append(orderItems, entities.OrderItem{
			ProductID:        cartItem.ProductID,
			Quantity:         cartItem.Quantity,
			SalePriceOnOrder: cartItem.FkProduct.FkDimensionalVariation.FkColourVariant.SalePrice,
		})
	}

	referenceNo, err := tools.MakeRandomUUID()
	if err != nil {
		fmt.Println("Error occured while generating reference number")
		return nil, nil, err
	}

	shippingCharges := getShippingCharge(address, totalProductValue)

	if req.CouponID != 0 {
		coupon, err := uc.orderRepo.GetCouponByID(req.CouponID)
		if err != nil {
			return nil, nil, err
		}

		couponUsageCount, err := uc.orderRepo.GetCouponUsageCount(req.UserID, req.CouponID)
		if err != nil {
			return nil, nil, err
		}

		couponDiscount, err = getCouponDiscount(coupon, totalProductValue, couponUsageCount)
		if err != nil {
			return nil, nil, err
		}
	}

	finalAmount := totalProductValue + shippingCharges - couponDiscount

	if strings.ToUpper(req.PaymentMethod) == "COD" {
		if !config.DeliveryConfig.CashOnDeliveryAvailable {
			return nil, nil, e.ErrCODNotAvailable
		}
		if totalProductValue > config.DeliveryConfig.MaxOrderAmountForCOD {
			return nil, nil, e.ErrOrderExceedsMaxAmountForCOD
		}
		status = "placed"
	} else {

		status = "payment pending"
		transactionID, err = services.CreateRazorpayOrder(finalAmount, referenceNo)
		if err != nil {
			return nil, nil, err
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

	if order.Status == "payment pending" {
		order.ID, err = uc.orderRepo.MakeOrder_UpdateStock_ClearCart(&order, &orderItems)
		if err != nil {
			fmt.Println("Error occured while placing order")
			return nil, nil, err
		}
	} else {
		order.ID, err = uc.orderRepo.MakeOrder(&order, &orderItems)
		if err != nil {
			fmt.Println("Error occured while placing order")
			return nil, nil, err
		}
	}
	//make order

	orderInfo.OrderDetails = order

	proceedToPaymentInfo := response.ProceedToPaymentInfo{
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

	return &orderInfo, &proceedToPaymentInfo, nil
}

// GetOrdersOfUser
func (uc *OrderUseCase) GetOrdersOfUser(userID uint, page int, limit int) (*[]response.ResponseOrderInfo, string, *e.Error) {
	var orders *[]entities.DetailedOrderInfo
	var responseOrders []response.ResponseOrderInfo
	offset := (page - 1) * limit
	orders, err := uc.orderRepo.GetOrdersOfUser(userID, offset, limit)
	if err != nil {
		fmt.Println("Error occured while getting orders")
		return &responseOrders, "Error occured while getting orders", err
	}

	err = copier.Copy(&responseOrders, &orders)
	if err != nil {
		return &responseOrders, "Error occured while copying orders to responseOrders", err
	}

	return &responseOrders, "Orders fetched successfully", nil
}

// GetAllOrders
func (uc *OrderUseCase) GetOrders(page int, limit int) (*[]response.ResponseOrderInfo, *e.Error) {
	var orders *[]entities.DetailedOrderInfo
	var responseOrders []response.ResponseOrderInfo
	var err error
	offset := (page - 1) * limit
	orders, err = uc.orderRepo.GetOrders(offset, limit)
	if err != nil {
		return &responseOrders, err
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
func (uc *OrderUseCase) CancelOrderByUser(orderID uint, userID uint) *e.Error {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return err
	}
	if !orderExists {
		return errors.New("order doesn't exist by ID")
	}

	//check if order belongs to userID
	userIDFromOrder, err := uc.orderRepo.GetUserIDByOrderID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting userID")
		return err
	}

	//check if userID in argument and userID from order match
	if userID != userIDFromOrder {
		return errors.New("order doesn't belong to user")
	}

	//check if order is already cancelled
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return err
	}

	if orderStatus != "placed" && orderStatus != "payment pending" {
		if orderStatus == "cancelled" {
			message := "Order is already in '" + orderStatus + "' status"
			return errors.New(message)
		} else {
			message := "Cannot cancel. Order is in '" + orderStatus + "' status"
			return errors.New(message)
		}
	}

	//cancel order, update stock, refund to wallet if already paid
	err = uc.orderRepo.CancelOrder(orderID)
	if err != nil {
		fmt.Println("Error occured while cancelling order")
		return err
	}

	return nil
}

// CancelOrderByAdmin(orderID uint) (string, error)
func (uc *OrderUseCase) CancelOrderByAdmin(orderID uint) *e.Error {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return err
	}
	if !orderExists {
		return errors.New("order doesn't exist by ID")
	}

	//check if order is already cancelled
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return err
	}
	if orderStatus != "placed" && orderStatus != "payment pending" {
		if orderStatus == "cancelled" {
			message := "Order is already in '" + orderStatus + "' status"
			return errors.New(message)
		} else {
			message := "Cannot cancel. Order is in '" + orderStatus + "' status"
			return errors.New(message)
		}
	}

	//cancel order, update stock, refund to wallet if already paid
	err = uc.orderRepo.CancelOrder(orderID)
	if err != nil {
		fmt.Println("Error occured while cancelling order")
		return err
	}

	return nil
}

// ReturnOrderRequestByUser
func (uc *OrderUseCase) ReturnOrderRequestByUser(orderID, userID uint) *e.Error {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return err
	}
	if !orderExists {
		return errors.New("order doesn't exist by ID")
	}

	//check if order belongs to userID
	userIDFromOrder, err := uc.orderRepo.GetUserIDByOrderID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting userID")
		return err
	}
	if userID != userIDFromOrder {
		return errors.New("order doesn't belong to user")
	}

	//check if order is in "delivered" status
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return err
	}
	if orderStatus != "delivered" {
		message := "Cannot return. Order is in '" + orderStatus + "' status"
		fmt.Println(message)
		return errors.New(message)
	}

	//return order
	err = uc.orderRepo.ReturnOrderRequest(orderID)
	if err != nil {
		fmt.Println("Error occured while returning order")
		return err
	}

	return nil
}

// MarkOrderAsReturned
func (uc *OrderUseCase) MarkOrderAsReturned(orderID uint) *e.Error {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return err
	}
	if !orderExists {
		return errors.New("order doesn't exist by ID")
	}

	//check if order is in "return requested" status
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return err
	}
	if orderStatus != "return requested" {
		message := "Cannot mark as returned. Order is in '" + orderStatus + "' status"
		fmt.Println(message)
		return errors.New(message)
	}

	//mark order as returned, update stock, refund to wallet
	err = uc.orderRepo.MarkOrderAsReturned(orderID)
	if err != nil {
		fmt.Println("Error occured while marking order as returned")
		return err
	}

	return nil
}

// MarkOrderAsDelivered
func (uc *OrderUseCase) MarkOrderAsDelivered(orderID uint) *e.Error {
	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return err
	}
	if !orderExists {
		return errors.New("order doesn't exist by ID")
	}

	//check if order is in "placed" status
	orderStatus, err := uc.orderRepo.GetOrderStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order status")
		return err
	}
	if orderStatus != "placed" {
		if orderStatus == "delivered" {
			message := "Order is already in '" + orderStatus + "' status"
			return errors.New(message)
		} else {
			message := "Cannot mark as delivered. Order is in '" + orderStatus + "' status"
			return errors.New(message)
		}
	}

	err = uc.orderRepo.MarkOrderAsDelivered(orderID)
	if err != nil {
		fmt.Println("Error occured while marking order as delivered")
		return err
	}

	return nil
}

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

func getCouponDiscount(coupon *entities.Coupon, orderValue float32, usageCount uint) (float32, *e.Error) {
	var discount float32
	switch {
	case coupon.Blocked:
		return 0, e.Error{Err: errors.New("coupon is blocked"), StatusCode: 403}
	case coupon.StartDate.After(time.Now()):
		return 0, e.Error{Err: errors.New("coupon is not active yet"), StatusCode: 403}
	case coupon.EndDate.Before(time.Now()):
		return 0, e.Error{Err: errors.New("coupon is expired"), StatusCode: 403}
	case coupon.MinOrderValue > orderValue:
		return 0, e.Error{Err: errors.New("order amount is less than required for coupon"), StatusCode: 403}
	case usageCount >= coupon.UsageLimit:
		return 0, e.Error{Err: errors.New("coupon usage limit exceeded"), StatusCode: 403}
	default:
		if coupon.Type == "fixed" {
			discount = myMath.RoundFloat32(coupon.MaxDiscount, 2)
		} else { //if coupon.Type == "percentage" {
			discount = max(orderValue*coupon.Percentage/100, coupon.MaxDiscount)
			discount = myMath.RoundFloat32(discount, 2)
		}
		return discount, nil
	}
}

// GetAddressForCheckout implements usecase.IOrderUC.
func (uc *OrderUseCase) GetAddressForCheckout(userID uint) (*[]entities.UserAddress, uint, float32, *e.Error) {
	quantity, totalValue, err := uc.cartRepo.GetQuantityAndPriceOfCart(userID)
	if err != nil {
		return nil, 0, 0, err
	}

	addresses, err := uc.userRepo.GetUserAddresses(userID)
	if err != nil {
		return nil, 0, 0, err
	}

	return addresses, quantity, totalValue, nil
}

// SetAddressGetCoupons implements usecase.IOrderUC.
func (uc *OrderUseCase) SetAddressGetCoupons(userID uint, req *request.SetAddressForCheckOutReq) (*response.SetAddrGetCouponsResponse, *e.Error) {
	// var resp response.SetAddrGetCouponsResponse
	address, err := uc.userRepo.GetUserAddress(req.AddressID)
	if err != nil {
		return nil, err
	}
	if address.UserID != userID {
		return nil, &e.Error {Err:errors.New("address doesn't belong to user"), StatusCode: 401}
	}

	totalQuantiy, totalProductsValue, err := uc.cartRepo.GetQuantityAndPriceOfCart(userID)
	if err != nil {
		return nil, err
	}

	//get shipping charge
	shippingCharge := getShippingCharge(address, totalProductsValue)
	if err != nil {
		return nil, err
	}

	//get coupons
	coupons, err := uc.orderRepo.GetActiveCoupons()
	if err != nil {
		return nil, err
	}
	var respCoupons []response.ResponseCoupon
	err = copier.Copy(&respCoupons, &coupons)
	if err != nil {
		fmt.Println("Error occured while copying coupons to responseCoupons, error:", err)
		return nil, err
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

	return &response, nil
}

// SetCouponGetPaymentMethods implements usecase.IOrderUC.
func (uc *OrderUseCase) SetCouponGetPaymentMethods(userID uint, req *request.SetCouponForCheckoutReq) (*response.GetPaymentMethodsForCheckoutResponse, *e.Error) {

	address, err := uc.userRepo.GetUserAddress(req.AddressID)
	if err != nil {
		return nil, err
	} else if address.UserID != userID {
		return nil, &e.Error {Err:errors.New("address doesn't belong to user"), StatusCode: 401}
	}

	//get coupon by id

	totalQuantiy, totalProductsValue, err := uc.cartRepo.GetQuantityAndPriceOfCart(userID)
	if err != nil {
		return nil, err
	}

	shippingCharge := getShippingCharge(address, totalProductsValue)
	if err != nil {
		return nil, err
	}

	coupon, err := uc.orderRepo.GetCouponByID(req.CouponID)
	if err != nil {
		return nil, err
	}

	usageCount, err := uc.orderRepo.GetCouponUsageCount(userID, req.CouponID)
	if err != nil {
		return nil, err
	}

	couponDiscount, err := getCouponDiscount(coupon, totalProductsValue, usageCount)
	if err != nil {
		return nil, err
	}

	walletBalance, err := uc.userRepo.GetWalletBalance(userID)
	if err != nil {
		return nil, err
	}

	grandTotal := totalProductsValue - couponDiscount + shippingCharge
	paymentMethods, codAvailability, codAvailabilityNote := getPaymentMethods(grandTotal)

	var respCoupon response.ResponseCoupon
	err = copier.Copy(&respCoupon, &coupon)
	if err != nil {
		fmt.Println("Error occured while copying coupon to responseCoupon")
		return nil, err
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

	return &resp, nil
}

func (uc *OrderUseCase) GetInvoiceOfOrder(userID uint, orderID uint) (*string, *e.Error) {

	//check if order exists
	orderExists, err := uc.orderRepo.DoOrderExistByID(orderID)
	if err != nil {
		fmt.Println("Error occured while checking if order exists")
		return nil, err
	}
	if !orderExists {
		return nil, &e.Error {Err:errors.New("order doesn't exist by ID"), StatusCode: 400}
	}

	//check if order belongs to userID
	userIDFromOrder, err := uc.orderRepo.GetUserIDByOrderID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting userID")
		return nil, err
	}

	if userID != userIDFromOrder {
		return nil, &e.Error {Err:errors.New("order doesn't belong to user"), StatusCode: 401}
	}

	//check if payment status is "paid"
	paymentStatus, err := uc.orderRepo.GetPaymentStatusByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting payment status")
		return nil, err
	}
	if paymentStatus != "paid" {
		message := "Cannot generate invoice. Payment status is '" + paymentStatus + "'"
		return nil, &e.Error {Err:errors.New(message), StatusCode: 400}
	}

	//get orderInfo
	order, err := uc.orderRepo.GetOrderSummaryByID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order summary")
		return nil, err
	}
	orderItems, err := uc.orderRepo.GetOrderItemsPQRByOrderID(orderID)
	if err != nil {
		fmt.Println("Error occured while getting order items")
		return nil, err
	}

	//get user details
	userInfo, err := uc.userRepo.GetUserBasicInfoByID(userID)
	if err != nil {
		fmt.Println("Error occured while getting user basic info")
		return nil, err
	}

	invoiceInfo := response.InvoiceInfo{
		OrderDetails: *order,
		OrderItems:   *orderItems,
		UserInfo:     *userInfo,
	}

	pdf := makeInvoicePDF(&invoiceInfo)

	if !config.ShouldUploadInvoice {
		fmt.Println("Uploading invoice to cloud is disabled. Invoice will be saved locally and link will be provided from that.")
		// Output the PDF to a file
		outputPath := filepath.Join(config.ExecutableDir, "testKit/invoiceOutput.pdf")
		err = pdf.OutputFileAndClose(outputPath)
		if err != nil {
			fmt.Println("Error saving PDF:", err)
			return nil, err
		} else {
			return &outputPath, nil
		}
	} else {
		tempFileName, err := tools.MakeRandomUUID()
		if err != nil {
			fmt.Println("Error occured while generating random UUID")
			return nil, err
		}
		tempFilePath := filepath.Join(os.TempDir(), tempFileName+"invoice.pdf")
		defer os.Remove(tempFilePath)
		err = pdf.OutputFileAndClose(tempFilePath)
		if err != nil {
			fmt.Println("Error saving PDF:", err)
			return nil, err
		}
		url, err := uc.orderRepo.UploadInvoice(tempFilePath, fmt.Sprint("invoice", orderID))
		if err != nil {
			fmt.Println("Error uploading PDF:", err)
			return nil, err
		}

		fmt.Println("url: ", url)

		return &url, nil
	}

}

func makeInvoicePDF(data *response.InvoiceInfo) *gofpdf.Fpdf {
	var billingUserStr, orderInfoStr, paymentInfoStr, shippingAddrStr string
	var productTotalStr, shippingChargeStr, couponDiscountStr, netSumStr string
	type item struct {
		Name      string
		MRP       string
		SalePrice string
		Quantity  string
		Net       string
	}

	orderItems := make([]item, len(data.OrderItems))

	{ //set strings/datas
		billingUserStr = "Name: " + data.UserInfo.FirstName + " " + data.UserInfo.LastName +
			"\nEmail: " + data.UserInfo.Email +
			"\nMobile: " + data.UserInfo.Phone

		orderInfoStr = "Ref No: " + data.OrderDetails.ReferenceNo +
			"\nOrder Date: " + data.OrderDetails.OrderDateAndTime.Format("02-01-2006") +
			"\nOrder Time: " + data.OrderDetails.OrderDateAndTime.Format("03:04PM")

		paymentInfoStr = "Payment Method: " + data.OrderDetails.PaymentMethod +
			"\nTransaction ID: " + data.OrderDetails.TransactionID

		shippingAddrStr = "Name: " + data.OrderDetails.FkAddress.FirstName + " " + data.OrderDetails.FkAddress.LastName +
			"\nEmail: " + data.OrderDetails.FkAddress.Email +
			"\nMobile: " + data.OrderDetails.FkAddress.Phone +
			"\n" + data.OrderDetails.FkAddress.Street +
			"\n" + data.OrderDetails.FkAddress.City +
			"\n" + data.OrderDetails.FkAddress.State +
			"\nPIN CODE- " + fmt.Sprint(data.OrderDetails.FkAddress.Pincode)

		for i, dataItem := range data.OrderItems {
			orderItems[i] = item{
				Name:      dataItem.ProductName,
				MRP:       fmt.Sprint(dataItem.MRP),
				SalePrice: fmt.Sprint(dataItem.SalePrice),
				Quantity:  fmt.Sprint(dataItem.Quantity),
				Net:       fmt.Sprint(dataItem.SalePrice * float32(dataItem.Quantity)),
			}
		}

		productTotalStr = fmt.Sprint(data.OrderDetails.OriginalAmount)
		shippingChargeStr = fmt.Sprint(data.OrderDetails.ShippingCharge)
		couponDiscountStr = fmt.Sprint(data.OrderDetails.CouponDiscount)
		netSumStr = fmt.Sprint(data.OrderDetails.FinalAmount)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 18)
	multiCellHeight := 7.0
	headLineSpacing := 8.0

	{ //left top
		leftWidth := 110.0
		{ // Logo
			logoPath := filepath.Join(config.ExecutableDir, "config/invoiceLogo.png")
			pdf.Image(logoPath, 10, 10, 70, 0, false, "", 0, "")
		}
		{ // Billing-to info
			pdf.SetY(35)
			pdf.SetFont("Arial", "B", 12)
			pdf.CellFormat(0, headLineSpacing, "Billing to:", "", 1, "L", false, 0, "")

			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(leftWidth, multiCellHeight, billingUserStr, "", "L", false)
		}
		{ // Order info
			pdf.SetY(pdf.GetY() + 7)
			pdf.SetFont("Arial", "B", 12)
			pdf.CellFormat(0, headLineSpacing, "Order Information:", "", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(leftWidth, multiCellHeight, orderInfoStr, "", "L", false)
		}
		{ // Payment info
			pdf.SetY(pdf.GetY() + 7)
			pdf.SetFont("Arial", "B", 12)
			pdf.CellFormat(0, headLineSpacing, "Payment Information:", "", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(leftWidth, multiCellHeight, paymentInfoStr, "", "L", false)
			pdf.Ln(10)
		}
	}
	{ //right top
		{ // Invoice label
			pdf.SetXY(85, 10)
			pdf.SetFont("Arial", "B", 18)
			pdf.CellFormat(0, 10, "Order Invoice #", "", 1, "R", false, 0, "")
			pdf.SetFont("Arial", "", 8)
			pdf.CellFormat(0, 10, "Order Ref No: "+data.OrderDetails.ReferenceNo, "B", 0, "R", false, 0, "")
		}

		{ //shipping address
			pdf.SetXY(125, pdf.GetY()+20)
			pdf.SetFont("Arial", "B", 12)
			pdf.CellFormat(80, headLineSpacing, "Shipping Address:", "", 2, "L", false, 0, "")
			pdf.SetX(pdf.GetX() + 3)
			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(70, multiCellHeight, shippingAddrStr, "", "L", false)
		}
	}

	{ // Item details table
		pdf.SetXY(5, 150)
		pdf.SetFillColor(170, 170, 170)
		pdf.CellFormat(88, 10, "Item", "B", 0, "L", true, 0, "")
		pdf.CellFormat(30, 10, "MRP", "B", 0, "L", true, 0, "")
		pdf.CellFormat(30, 10, "Sale Price", "B", 0, "L", true, 0, "")
		pdf.CellFormat(20, 10, "Quantity", "B", 0, "L", true, 0, "")
		pdf.CellFormat(30, 10, "Net", "B", 1, "L", true, 0, "")

		for _, item := range orderItems {
			pdf.SetX(5)
			// pdf.CellFormat(88, 10, item.Name, "", 0, "L", false, 0, "")
			pdf.MultiCell(88, 10, item.Name, "", "L", false)
			pdf.SetXY(93, pdf.GetY()-10)
			pdf.CellFormat(28, 10, item.MRP, "", 0, "R", false, 0, "")
			pdf.CellFormat(30, 10, item.SalePrice, "", 0, "R", false, 0, "")
			pdf.CellFormat(17, 10, item.Quantity, "", 0, "R", false, 0, "")
			pdf.CellFormat(33, 10, item.Net, "", 1, "R", false, 0, "")
		}

		pdf.Line(5, pdf.GetY(), 200, pdf.GetY())
	}

	{ //Add additional charges
		widthTitle := 40.0
		widthValue := 30.0
		lineSpacing := 8.0
		x := 131.0
		pdf.SetX(x)
		pdf.CellFormat(widthTitle, lineSpacing, "Products Total:", "", 0, "L", false, 0, "")
		pdf.CellFormat(widthValue, lineSpacing, productTotalStr, "", 1, "R", false, 0, "")

		pdf.SetX(x)
		pdf.CellFormat(widthTitle, lineSpacing, "Shipping Charge:", "", 0, "L", false, 0, "")
		pdf.CellFormat(widthValue, lineSpacing, shippingChargeStr, "", 1, "R", false, 0, "")

		pdf.SetX(x)
		pdf.CellFormat(widthTitle, lineSpacing, "Coupon Discount:", "", 0, "L", false, 0, "")
		pdf.CellFormat(widthValue, lineSpacing, couponDiscountStr, "", 1, "R", false, 0, "")

		pdf.SetX(x)
		pdf.SetFontStyle("B")
		pdf.CellFormat(widthTitle, lineSpacing, "Net Sum:", "T", 0, "L", false, 0, "")
		pdf.CellFormat(widthValue, lineSpacing, netSumStr, "T", 1, "R", false, 0, "")
	}

	return pdf
}
