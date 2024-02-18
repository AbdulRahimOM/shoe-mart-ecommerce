package services

import (
	"MyShoo/internal/config"
	"fmt"

	razorpay "github.com/razorpay/razorpay-go"
	razorpayUtils "github.com/razorpay/razorpay-go/utils"
)

func CreateRazorpayOrder(amount float32, referenceNo string) (string, error) {
	client := razorpay.NewClient(config.RazorpayKeyId, config.RazorpayKeySecret)
	data := map[string]interface{}{
		"amount":   uint(amount * 100),
		"currency": "INR",
		// "receipt":         "some_receipt_id",	//may add at future, along with invoice
		"partial_payment": false,
		"notes": map[string]interface{}{
			"ReferenceNo": referenceNo,
		},
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		fmt.Println("error creating razorpay order. err= ", err)
		return "", err
	}
	transactionID := body["id"].(string)
	return transactionID, nil

}

func VerifyPayment(razorpay_order_id string, razorpay_payment_id string, signature string) bool {
	params := map[string]interface{}{
		"razorpay_order_id":   razorpay_order_id,
		"razorpay_payment_id": razorpay_payment_id,
	}

	isValid := razorpayUtils.VerifyPaymentSignature(params, signature, config.RazorpayKeySecret)
	if isValid {
		fmt.Println("Payment is valid")
		return true
	} else {
		fmt.Println("Payment invalid")
		return false
	}

}
