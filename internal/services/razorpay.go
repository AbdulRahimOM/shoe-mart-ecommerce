package services

import (
	"MyShoo/internal/domain/entities"
	"fmt"
	"os"

	razorpay "github.com/razorpay/razorpay-go"
)

func CreateRazorpayOrder(order entities.Order) (string, error) {

	var razorpayID = os.Getenv("RAZORPAY_KEY_ID")
	var razorpaySecret = os.Getenv("RAZORPAY_KEY_SECRET")
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

func VerifyPayment() {

}
