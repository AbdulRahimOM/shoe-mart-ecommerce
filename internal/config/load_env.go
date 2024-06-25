package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)
var Port string
var IsLocalHostMode bool = false
var ExecutableDir string

// var DB_URL string
var DbURL string

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

// Development mode.......................
var ShouldUploadExcel bool
var ShouldRenderPaymentPage bool
var ShouldUploadInvoice bool

var relativeEnvPath  string ="config/envs/.env"

func init() {

}

func LoadEnvVariables() error {
	var err error

	ExecutableDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("error getting current directory path: %v", err)
	}

	//try to load .env file in 'execution of binary' mode
	errBinaryExecMode := godotenv.Load(filepath.Join(ExecutableDir, relativeEnvPath)) //env file is presumed to be alongside the executable
	if errBinaryExecMode != nil {
		goto retryWithGoRunMode
	}else{
		initiateEnvValues()
		return nil
	}

retryWithGoRunMode:
	//try to load .env file in 'go run' mode
	if errGoRunMode := godotenv.Load(relativeEnvPath); errGoRunMode != nil {
		return fmt.Errorf("error loading .env file by either modes. \nerr from binary mode: %v\n, err from go run mode: %v", errBinaryExecMode, errGoRunMode)
	} else {
		ExecutableDir = ""
		initiateEnvValues()
		return nil
	}
}

func initiateEnvValues() {
	// Development mode.......................
	if os.Getenv("LOCAL_HOST_MODE") == "true" {
		IsLocalHostMode = true
	}

	Port = os.Getenv("PORT")

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

	if os.Getenv("UPLOAD_EXCEL") == "true" {
		ShouldUploadExcel = true
	}
	if os.Getenv("RENDER_PAYMENT_PAGE") == "true" {
		ShouldRenderPaymentPage = true
	}
	if os.Getenv("UPLOAD_INVOICE") == "true" {
		ShouldUploadInvoice = true
	}
}
