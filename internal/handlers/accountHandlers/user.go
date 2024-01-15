package handlers

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwttoken "MyShoo/pkg/jwt_tokens"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase usecaseInterface.IUserUC
}

func NewUserHandler(useCase usecaseInterface.IUserUC) *UserHandler {
	return &UserHandler{UserUseCase: useCase}
}

// definefunctions those which you  have called in user routes
func (h *UserHandler) GetLogin(c *gin.Context) {
	fmt.Println("Handler ::: GET login handler")

	c.JSON(http.StatusOK, "token")
}
func (h *UserHandler) GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hai":   "dfdf",
		"hello": "hello",
	})
}

func (h *UserHandler) PostSignUp(c *gin.Context) {
	fmt.Println("=============\nentered POST sign-up handler")

	var signUpReq requestModels.UserSignUpReq

	if err := c.Bind(&signUpReq); err != nil {
		fmt.Println("\n", "Error occured while signing up. Error while binding request"+err.Error(), "\n.")

		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error while binding request",
			Error:   err.Error(),
		})

		return
	}

	//validation
	if err := requestValidation.ValidateRequest(signUpReq); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#", //how to inform that this particular field is wrong, like email is wrong, i have to find a way to do that
			Error:   errResponse,
			Token:   "",
		})

		return
	}

	token, err := h.UserUseCase.SignUp(&signUpReq)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "Error occured while signing up. Try again. Error:" + err.Error() ////////////////////////////////
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})

		return
	} else {
		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "",
			Token:   *token,
		})
	}
}

// -----------------------------------------------------------------------------------------------------------------
func (h *UserHandler) PostLogIn(c *gin.Context) {
	fmt.Println("=============\nentered \"POST login\" handler")

	var signInReq requestModels.UserSignInReq

	if err := c.ShouldBindJSON(&signInReq); err != nil {
		fmt.Println("\nerror binding the requewst\n.")
		errResponse := "error binding the requewst. Try again. Error:" + err.Error()
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})

		return
	}

	//validation
	if err := requestValidation.ValidateRequest(signInReq); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})

		return
	}

	token, err := h.UserUseCase.SignIn(&signInReq)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "error while signing in"
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})

		return
	} else {
		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "",
			Token:   *token,
		})
	}
}

// -----------------------------------------------------------------------------------------------------------------
func (h *UserHandler) SendOtp(c *gin.Context) {
	fmt.Println("=============entered \"Send OTP \" handler======")

	user, ok := c.Get("UserModel")
	if !ok {
		fmt.Println("error getting user model from context")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error happened. Please login again",
			Error:   "Error getting user model from context",
			Token:   "",
		})
		return
	}
	userMap := user.(map[string]interface{})
	phone := userMap["phone"].(string)
	err := h.UserUseCase.SendOtp(phone)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "error while sending otp",
			Error:   err.Error(),
			Token:   "",
		})
		return
	} else {
		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "OTP sent successfully. Please verify",
		})
	}
}

// -----------------------------------------------------------------------------------------------------------------
func (h *UserHandler) VerifyOtp(c *gin.Context) {
	fmt.Println("=============entered \"Verify OTP \" handler======")
	var otpStruct struct {
		OTP string `json:"otp" validate:"required,number"`
	}
	if err := c.Bind(&otpStruct); err != nil {
		fmt.Println("\nerror binding the requewst\n.")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   "error binding the request. Error:" + err.Error(),
		})
		return
	}
	//validation
	if err := requestValidation.ValidateRequest(otpStruct); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})
		return
	}

	user, ok := c.Get("UserModel")
	if !ok {
		fmt.Println("error getting user model from context")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error happened. Please try again",
			Error:   "Error getting user model from context",
		})
		return
	}
	phone := user.(map[string]interface{})["phone"].(string)
	email := user.(map[string]interface{})["email"].(string)
	isVerified, err := h.UserUseCase.VerifyOtp(phone, email, otpStruct.OTP)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "error occured while verifying otp. Please try again",
			Error:   err.Error(),
		})
		return
	}
	if isVerified {

		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "OTP verified successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "OTP verification failed. Please try again",
		})
	}
}

// Add address
func (h *UserHandler) AddUserAddress(c *gin.Context) {
	fmt.Println("Handler ::: add address handler")
	var req requestModels.AddUserAddress

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//check if userID in token and request body match
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding to cart. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if userID != req.UserID {
		fmt.Println("User ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "User ID in token and request body do not match",
		})
		return
	}

	if err := h.UserUseCase.AddUserAddress(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding address. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Address added successfully",
	})

}

// Edit address
func (h *UserHandler) EditUserAddress(c *gin.Context) {
	fmt.Println("Handler ::: edit address handler")
	var req requestModels.EditUserAddress

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//check if userID in token and request body match
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing address. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if userID != req.UserID {
		fmt.Println("User ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "User ID in token and request body do not match",
		})
		return
	}

	if err := h.UserUseCase.EditUserAddress(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing address. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Address edited successfully",
	})
}

// DeleteUserAddress
func (h *UserHandler) DeleteUserAddress(c *gin.Context) {
	fmt.Println("Handler ::: delete address handler")
	var req requestModels.DeleteUserAddress

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//check if userID in token and request body match
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error deleting address. Try Again",
			Error:   err.Error(),
		})
		return
	}
	if userID != req.UserID {
		fmt.Println("User ID in token and request body do not match. Corrupted request!!")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Corrupted request. Try Again",
			Error:   "User ID in token and request body do not match",
		})
		return
	}

	if err := h.UserUseCase.DeleteUserAddress(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error deleting address. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Address deleted successfully",
	})
}

// Get user addresses
func (h *UserHandler) GetUserAddresses(c *gin.Context) {
	fmt.Println("Handler ::: get user addresses handler")

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetUserAddressesResponse{
			Status:  "failed",
			Message: "Error getting addresses. Try Again",
			Error:   err.Error(),
		})
		return
	}

	addresses, err := h.UserUseCase.GetUserAddresses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetUserAddressesResponse{
			Status:  "failed",
			Message: "Error getting addresses. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetUserAddressesResponse{
		Status:    "success",
		Message:   "Addresses fetched successfully",
		Addresses: *addresses,
	})
}

// GetProfile
func (h *UserHandler) GetProfile(c *gin.Context) {
	fmt.Println("Handler ::: get profile handler")

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetProfileResponse{
			Status:  "failed",
			Message: "Error getting profile. Try Again",
			Error:   err.Error(),
		})
		return
	}

	profile, err := h.UserUseCase.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetProfileResponse{
			Status:  "failed",
			Message: "Error getting profile. Try Again",
			Error:   err.Error(),
		})
		return
	}

	var addresses *[]entities.UserAddress
	addresses, err = h.UserUseCase.GetUserAddresses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetProfileResponse{
			Status:  "failed",
			Message: "Error getting addresses. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetProfileResponse{
		Status:  "success",
		Message: "Profile fetched successfully",
		Profile: struct {
			UserDetails entities.UserDetails   `json:"userDetails"`
			Addresses   []entities.UserAddress `json:"addresses"`
		}{
			UserDetails: *profile,
			Addresses:   *addresses,
		},
	})
}

// EditProfile
func (h *UserHandler) EditProfile(c *gin.Context) {
	fmt.Println("Handler ::: edit profile handler")
	var req requestModels.EditProfileReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing profile. Try Again",
			Error:   err.Error(),
		})
		return
	}

	if err := h.UserUseCase.EditProfile(userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error editing profile. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Profile edited successfully",
	})
}

// GetResetPassword
func (h *UserHandler) SendOtpForPWChange(c *gin.Context) {
	fmt.Println("Handler ::: get reset password handler")

	var req requestModels.ApplyForPasswordResetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//get user info using email
	var user *entities.User
	user, err := h.UserUseCase.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SMET{
			Status:  "failed",
			Message: "Error getting reset password. Try Again",
			Error:   err.Error(),
		})
		return
	}

	token, err := h.UserUseCase.SendOtpForPWChange(user)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "error while sending otp",
			Error:   err.Error(),
		})
		return
	} else {

		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "OTP sent successfully. Please verify",
			Token:   *token,
		})
	}
}

func (h *UserHandler) VerifyOtpForPWChange(c *gin.Context) {
	fmt.Println("=============entered \"Verify OTP for pw change \" handler======")

	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString) //
	secretKey := os.Getenv("SECRET_KEY")
	isTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, secretKey)
	if !isTokenValid {
		fmt.Println("token is invalid")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}
	//getting claims
	claims, ok := tokenClaims.(*jwttoken.CustomClaims)
	if !ok {
		fmt.Println("claims type assertion failed")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	//checking if role is user
	if claims.Role != "user" {
		fmt.Println("role is not user")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}
	status := claims.Model.(map[string]interface{})["Status"].(string)
	if status != "PW change requested, otp not verified" {
		fmt.Println("status is not PW change requested, otp not verified")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	phone := claims.Model.(map[string]interface{})["Phone"].(string)
	id := uint(claims.Model.(map[string]interface{})["ID"].(float64))

	var otpStruct struct {
		OTP string `json:"otp" validate:"required,number"`
	}
	if err := c.Bind(&otpStruct); err != nil {
		fmt.Println("\nerror binding the requewst\n.")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   "error binding the request. Error:" + err.Error(),
		})
		return
	}
	//validation
	if err := requestValidation.ValidateRequest(otpStruct); err != nil {
		fmt.Println("\n\nerror validating the request\n.")
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
			Token:   "",
		})
		return
	}

	isVerified, newtoken, err := h.UserUseCase.VerifyOtpForPWChange(id, phone, otpStruct.OTP)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "error occured while verifying otp. Please try again",
			Error:   err.Error(),
		})
		return
	}
	if isVerified {
		fmt.Println("New token: ", *newtoken)
		c.JSON(http.StatusOK, response.SMET{
			Status:  "success",
			Message: "OTP verified successfully",
			Token:   *newtoken,
		})
	} else {
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "OTP verification failed. Please try again",
		})
	}
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	fmt.Println("before entering reset password handler")
	// myshoo.Test()
	fmt.Println("Handler ::: reset password handler")

	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString)
	secretKey := os.Getenv("SECRET_KEY")
	isTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, secretKey)
	if !isTokenValid {
		fmt.Println("token is invalid!")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}
	fmt.Println("************2")
	//getting claims
	claims, ok := tokenClaims.(*jwttoken.CustomClaims)
	if !ok {
		fmt.Println("claims type assertion failed")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	//checking if role is user
	if claims.Role != "user" {
		fmt.Println("role is not user")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}
	fmt.Println("************3")
	status := claims.Model.(map[string]interface{})["Status"].(string)
	if status != "PW change requested, otp verified" {
		fmt.Println("status is not PW change requested, otp not verified")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	var req requestModels.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error binding request. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   errResponse,
		})
		return
	}

	//change password
	id := uint(claims.Model.(map[string]interface{})["ID"].(float64))
	if err := h.UserUseCase.ResetPassword(id, &req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, response.SMET{
			Status:  "failed",
			Message: "Error resetting password. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SMET{
		Status:  "success",
		Message: "Password reset successfully",
	})
}
