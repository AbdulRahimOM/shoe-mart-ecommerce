package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var LocalHostMode string
var ExecutableDir string

// var err error
var envPath string

// var DB_URL string
var DbURL  string

// JWT token generation.....................
var SecretKey string

// Twilio OTP generation.....................
var TwilioAccountSid string
var TwilioAuthToken string
var TwilioServiceSid string

// Cloudinary...............................
var CloudinaryCloudName string
var CloudinaryApiKey string
var CloudinaryApiSecret string

// Razorpay...............................
var RazorpayKeyId string
var RazorpayKeySecret string

//Development mode.......................
var UploadExcel string
var RenderPaymentPage string
var UploadInvoice string


func init() {
	var err error
	if LocalHostMode == "true" {
		ExecutableDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			panic("Couldn't get current(executable) directory")
		}
	}
}

func LoadEnvVariables() error {
	envPath = filepath.Join(ExecutableDir, ".env") //env file is presumed to be alongside the executable
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Couldn't load env variables")
		return err
	}

	initiateEnvValues()
	return nil
}

func initiateEnvValues() {

	// Database URL.....................
	DbURL = os.Getenv("DB_URL")

	// JWT token generation.....................
	SecretKey = os.Getenv("SECRET_KEY")

	// Twilio OTP generation.....................
	TwilioAccountSid = os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	TwilioServiceSid = os.Getenv("TWILIO_SERVICE_SID")

	// Cloudinary...............................
	CloudinaryCloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
	CloudinaryApiKey = os.Getenv("CLOUDINARY_API_KEY")
	CloudinaryApiSecret = os.Getenv("CLOUDINARY_API_SECRET")

	// Razorpay...............................
	RazorpayKeyId = os.Getenv("RAZORPAY_KEY_ID")
	RazorpayKeySecret = os.Getenv("RAZORPAY_KEY_SECRET")

	UploadExcel = os.Getenv("UPLOAD_EXCEL")
	RenderPaymentPage = os.Getenv("RENDER_PAYMENT_PAGE")
	UploadInvoice = os.Getenv("UploadInvoice")
}
