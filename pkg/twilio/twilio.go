package otpManager

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func SendOtp(phone string) error {
	var serviceSID = os.Getenv("TWILIO_SERVICE_SID")
	var accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
	var authToken = os.Getenv("TWILIO_AUTH_TOKEN")
	// client := twilio.NewRestClient(twilio.ClientParams{Username: accountSid, Password: authToken})
	fmt.Println("**********entered send otp*****")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSid, Password: authToken})
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	fmt.Println("phone=", phone)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(serviceSID, params)
	if err != nil {
		fmt.Println(err.Error())
		// fmt.Println("**********middle3*****")
		// fmt.Println("accountSid=", accountSid)
		// fmt.Println("authToken=", authToken)
		// fmt.Println("serviceSID=", serviceSID)

		return err
	} else {
		if resp.Status != nil {
			fmt.Println("hh", *resp.Status)
		} else {
			fmt.Println("kk", resp.Status)
		}
		return nil
	}
}

func VerifyOtp(phone string, otp string) (bool, error) {
	var serviceSID = os.Getenv("TWILIO_SERVICE_SID")
	var accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
	var authToken = os.Getenv("TWILIO_AUTH_TOKEN")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSid, Password: authToken})
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(otp)

	resp, err := client.VerifyV2.CreateVerificationCheck(serviceSID, params)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	} else {
		if resp.Status != nil {
			fmt.Println(*resp.Status)
		} else {
			fmt.Println(resp.Status)
		}
		return true, nil
	}
}
