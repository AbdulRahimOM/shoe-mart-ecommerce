package accounthandler

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecase "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwttoken "MyShoo/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase usecase.IUserUC
}

func NewUserHandler(useCase usecase.IUserUC) *UserHandler {
	return &UserHandler{UserUseCase: useCase}
}

// to get user login page
// @Summary Get user login page
// @Description Get user login page
// @Tags User/Session/Login
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} string
// @Router /login [get]
func (h *UserHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

// @Summary User Sign Up Handler
// @Description User Sign Up Handler
// @Tags User/Session/SignUp
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.UserSignUpReq{} true "User Sign Up Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /signup [post]
func (h *UserHandler) PostSignUp(c *gin.Context) {

	var signUpReq request.UserSignUpReq
	if err := c.Bind(&signUpReq); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(signUpReq); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}
	token, err := h.UserUseCase.SignUp(&signUpReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("", err))
		return
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Status:  "success",
			Message: "",
			Token:   *token,
		})
	}
}

// @Summary User Sign In Handler
// @Description User Sign In Handler
// @Tags User/Session/Login
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.UserSignInReq{} true "User Sign In Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /login [post]
func (h *UserHandler) PostLogIn(c *gin.Context) {

	var signInReq request.UserSignInReq
	if err := c.ShouldBindJSON(&signInReq); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(signInReq); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	token, err := h.UserUseCase.SignIn(&signInReq)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "error while signing in"
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})

		return
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Status:  "success",
			Message: "",
			Token:   *token,
		})
	}
}

// @Summary Send OTP
// @Description Send OTP
// @Tags User/Session/SignUp
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /sendotp [get]
func (h *UserHandler) SendOtp(c *gin.Context) {

	user, ok := c.Get("UserModel")
	if !ok {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error happened. Please login again",
			Error:   "Error getting user model from context",
		})
		return
	}
	userMap := user.(map[string]interface{})
	phone := userMap["phone"].(string)
	err := h.UserUseCase.SendOtp(phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("", err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("OTP sent successfully"))
	}
}

// @Summary Verify OTP
// @Description Verify OTP
// @Tags User/Session/SignUp
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.VerifyOTPReq{} true "Verify OTP Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /verifyotp [post]
func (h *UserHandler) VerifyOtp(c *gin.Context) {

	var otpStruct request.VerifyOTPReq
	if err := c.Bind(&otpStruct); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(otpStruct); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	user, ok := c.Get("UserModel")
	if !ok {
		fmt.Println("error getting user model from context")
		c.JSON(http.StatusBadRequest, response.SME{
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
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "error occured while verifying otp. Please try again",
			Error:   err.Error(),
		})
		return
	}
	if isVerified {

		c.JSON(http.StatusOK, response.SM{
			Status:  "success",
			Message: "OTP verified successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, response.SM{
			Status:  "failed",
			Message: "OTP verification failed. Please try again",
		})
	}
}

// Add address
// @Summary Add address
// @Description Add address
// @Tags User/Address
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.AddUserAddress{} true "Add Address Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /addaddress [post]
func (h *UserHandler) AddUserAddress(c *gin.Context) {

	var req request.AddUserAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
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
// @Summary Edit address
// @Description Edit address
// @Tags User/Address
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.EditUserAddress{} true "Edit Address Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /editaddress [patch]
func (h *UserHandler) EditUserAddress(c *gin.Context) {

	var req request.EditUserAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
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
// @Summary Delete address
// @Description Delete address
// @Tags User/Address
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.DeleteUserAddress{} true "Delete Address Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /deleteaddress [delete]
func (h *UserHandler) DeleteUserAddress(c *gin.Context) {

	var req request.DeleteUserAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
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

	c.JSON(http.StatusOK, response.SM{
		Status:  "success",
		Message: "Address deleted successfully",
	})
}

// Get user addresses
// @Summary Get user addresses
// @Description Get user addresses
// @Tags User/Address
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetUserAddressesResponse{}
// @Failure 400 {object} response.SME{}
// @Router /addresses [get]
func (h *UserHandler) GetUserAddresses(c *gin.Context) {

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
// @Summary Get profile
// @Description Get profile
// @Tags User/Profile
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetProfileResponse{}
// @Failure 400 {object} response.SME{}
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {

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
// @Summary Edit profile
// @Description Edit profile
// @Tags User/Profile
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.EditProfileReq{} true "Edit Profile Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /editprofile [patch]
func (h *UserHandler) EditProfile(c *gin.Context) {

	var req request.EditProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
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
// @Summary Get reset password
// @Description Get reset password
// @Tags User/Session/Reset_password
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.ApplyForPasswordResetReq{} true "Apply For Password Reset Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /resetpassword [get]
func (h *UserHandler) SendOtpForPWChange(c *gin.Context) {

	var req request.ApplyForPasswordResetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//get user info using email
	var user *entities.User
	user, err := h.UserUseCase.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting reset password. Try Again",
			Error:   err.Error(),
		})
		return
	}

	token, err := h.UserUseCase.SendOtpForPWChange(user)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "error while sending otp",
			Error:   err.Error(),
		})
		return
	} else {

		c.JSON(http.StatusOK, response.SMT{
			Status:  "success",
			Message: "OTP sent successfully. Please verify",
			Token:   *token,
		})
	}
}

// VerifyOtpForPWChange
// @Summary Verify OTP for password change
// @Description Verify OTP for password change
// @Tags User/Session/Reset_password
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.VerifyOTPReq{} true "Verify OTP Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /resetpasswordverifyotp [post]
func (h *UserHandler) VerifyOtpForPWChange(c *gin.Context) {

	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString) //
	secretKey := os.Getenv("SECRET_KEY")
	isTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, secretKey)
	if !isTokenValid {
		fmt.Println("token is invalid")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}
	//getting claims
	claims, ok := tokenClaims.(*jwttoken.CustomClaims)
	if !ok {
		fmt.Println("claims type assertion failed")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}

	//checking if role is user
	if claims.Role != "user" {
		fmt.Println("role is not user")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}
	status := claims.Model.(map[string]interface{})["Status"].(string)
	if status != "PW change requested, otp not verified" {
		fmt.Println("status is not PW change requested, otp not verified")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}

	phone := claims.Model.(map[string]interface{})["Phone"].(string)
	id := uint(claims.Model.(map[string]interface{})["ID"].(float64))

	var otpStruct request.VerifyOTPReq
	if err := c.Bind(&otpStruct); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(otpStruct); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	isVerified, newtoken, err := h.UserUseCase.VerifyOtpForPWChange(id, phone, otpStruct.OTP)
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "error occured while verifying otp. Please try again",
			Error:   err.Error(),
		})
		return
	}
	if isVerified {
		fmt.Println("New token: ", *newtoken)
		c.JSON(http.StatusOK, response.SMT{
			Status:  "success",
			Message: "OTP verified successfully",
			Token:   *newtoken,
		})
	} else {
		c.JSON(http.StatusBadRequest, response.SM{
			Status:  "failed",
			Message: "OTP verification failed. Please try again",
		})
	}
}

// ResetPassword
// @Summary Reset password
// @Description User can provide new password after verifying OTP
// @Tags User/Session/Reset_password
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param req body req.ResetPasswordReq{} true "Reset Password Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /resetpassword [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {

	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString)
	secretKey := os.Getenv("SECRET_KEY")
	isTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, secretKey)
	if !isTokenValid {
		fmt.Println("token is invalid!")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}

	//getting claims
	claims, ok := tokenClaims.(*jwttoken.CustomClaims)
	if !ok {
		fmt.Println("claims type assertion failed")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}

	//checking if role is user
	if claims.Role != "user" {
		fmt.Println("role is not user")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		return
	}

	status := claims.Model.(map[string]interface{})["Status"].(string)
	if status != "PW change requested, otp verified" {
		fmt.Println("status is not PW change requested, otp not verified")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	var req request.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(err.Error(), e.ErrOnBindingReq))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME(fmt.Sprint(err), e.ErrOnValidation))
		return
	}

	//change password
	id := uint(claims.Model.(map[string]interface{})["ID"].(float64))
	if err := h.UserUseCase.ResetPassword(id, &req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error resetting password. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SM{
		Status:  "success",
		Message: "Password reset successfully",
	})
}
