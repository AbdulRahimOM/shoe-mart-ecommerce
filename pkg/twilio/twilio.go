package otpManager

import (
	"fmt"

	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/config"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func SendOtp(phone string) error {
	if config.TakeAllOtpAsValid {
		return nil
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: config.TwilioAccountSid, Password: config.TwilioAuthToken})
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(config.TwilioServiceSid, params)
	if err != nil {
		return err
	} else {
		if resp.Status != nil { //?
			fmt.Println("Otp sending status:", *resp.Status)
		} else {
			fmt.Println("kk", resp.Status)
		}
		return nil
	}
}

func VerifyOtp(phone string, otp string) (bool, error) {
	if config.TakeAllOtpAsValid {
		return true, nil
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: config.TwilioAccountSid, Password: config.TwilioAuthToken})
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(otp)
	fmt.Println("otp: ", otp)
	resp, err := client.VerifyV2.CreateVerificationCheck(config.TwilioServiceSid, params)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	} else {
		if resp.Status != nil {
			fmt.Println("*resp.Status: ", *resp.Status)
			if *resp.Status == "approved" {
				return true, nil
			}
		} else {
			fmt.Println("resp.Status: ", resp.Status)
		}
	}
	return false, err
}
