package services

import (
	"MyShoo/internal/domain/entities"
	"fmt"
	"os"

	razorpay "github.com/razorpay/razorpay-go"
	razorpayUtils "github.com/razorpay/razorpay-go/utils"
)

func CreateRazorpayOrder(order entities.Order) (string, error) {

	razorpayID := os.Getenv("RAZORPAY_KEY_ID")
	razorpaySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	client := razorpay.NewClient(razorpayID, razorpaySecret)
	data := map[string]interface{}{
		"amount":   uint(order.FinalAmount * 100),
		"currency": "INR",
		// "receipt":         "some_receipt_id",	//may add at future, along with invoice
		"partial_payment": false,
		"notes": map[string]interface{}{
			"ReferenceNo": order.ReferenceNo,
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
	razorpaySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	params := map[string]interface{}{
		"razorpay_order_id":   razorpay_order_id,
		"razorpay_payment_id": razorpay_payment_id,
	}

	secret := razorpaySecret
	isValid := razorpayUtils.VerifyPaymentSignature(params, signature, secret)
	if isValid {
		fmt.Println("Payment is valid")
		return true
	} else {
		fmt.Println("Payment invalid")
		return false
	}

}

