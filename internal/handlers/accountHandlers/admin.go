package handlers

import (
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminUseCase usecaseInterface.IAdminUC
}

func NewAdminHandler(useCase usecaseInterface.IAdminUC) *AdminHandler {
	return &AdminHandler{AdminUseCase: useCase}
}

// to get admin login page
// @Summary Get admin login page
// @Description Get admin login page
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /admin/login [get]
func (h *AdminHandler) GetAdminLogin(c *gin.Context) {
	fmt.Println("Handler ::: GET login handler")

	c.JSON(http.StatusOK, "This is admin login page. Enter credentials to login")
}

// to get admin home page : But not developed yet
// NOTE: this is just a sample
func (h *AdminHandler) GetAdminHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hai":   "dfdf",
		"hello": "hello admin",
	})
}

// to login admin
// @Summary Login admin
// @Description Login admin
// @Tags admin
// @Accept json
// @Produce json
// @Param adminSignInReq body requestModels.AdminSignInReq true "Admin Sign In Request"
// @Success 200 {object} response.SMET
// @Failure 400 {object} response.SMET
// @Router /admin/login [post]
func (h *AdminHandler) PostLogIn(c *gin.Context) {
	fmt.Println("=============\nentered \"POST login\" handler")

	var signInReq requestModels.AdminSignInReq
	fmt.Println("c.request.body=", c.Request.Body)
	if err := c.ShouldBindJSON(&signInReq); err != nil {
		fmt.Println("\nerror binding the requewst\n.")
		errResponse := "error binding the requewst. Try again. Error:" + err.Error()
		c.JSON(http.StatusBadRequest, response.SMET{
			Status:  "failed",
			Message: "Error occured. Please try again.",
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

	token, err := h.AdminUseCase.SignIn(&signInReq)
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

// to get users list
// @Summary Get users list
// @Description Get users list
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} response.GetUsersListResponse
// @Failure 400 {object} response.GetUsersListResponse
// @Router /admin/userslist [get]
func (h *AdminHandler) GetUsersList(c *gin.Context) {
	fmt.Println("Handler ::: \"GET users list\" handler")
	usersList, err := h.AdminUseCase.GetUsersList()
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "error while getting users list"
		c.JSON(http.StatusBadRequest, response.GetUsersListResponse{
			Status:    "failed",
			Message:   "Error occured while getting users list. Please try again.",
			Error:     errResponse,
			UsersList: nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.GetUsersListResponse{
			Status:    "success",
			Message:   "The list of users",
			Error:     "",
			UsersList: *usersList,
		})
	}
}

// to get sellers list
// @Summary Get sellers list
// @Description Get sellers list
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} response.GetSellersListResponse
// @Failure 400 {object} response.GetSellersListResponse
// @Router /admin/sellerslist [get]
func (h *AdminHandler) GetSellersList(c *gin.Context) {
	fmt.Println("Handler ::: \"GET sellers list\" handler")
	sellersList, err := h.AdminUseCase.GetSellersList()
	if err != nil {
		fmt.Println("\n\nHandler: error recieved from usecase\n\n.")
		errResponse := "error while getting sellers list"
		c.JSON(http.StatusBadRequest, response.GetSellersListResponse{
			Status:      "failed",
			Message:     "Error occured while getting sellers list. Please try again.",
			Error:       errResponse,
			SellersList: nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.GetSellersListResponse{
			Status:      "success",
			Message:     "The list of sellers",
			Error:       "",
			SellersList: *sellersList,
		})
	}
}

// to block a user
// @Summary Block user
// @Description Block user
// @Tags admin
// @Accept json
// @Produce json
// @Param unblockUserRequest body requestModels.BlockUserReq true "user"
// @Success 200 {object} response.BlockUserResponse
// @Failure 400 {object} response.BlockUserResponse
// @Router /admin/blockuser [post]
func (h *AdminHandler) BlockUser(c *gin.Context) {
	fmt.Println("Handler ::: \"Block user\" handler")

	//get user info from request
	var req requestModels.BlockUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse := "error binding the request. Try again. Error:" + err.Error()
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.BlockUserResponse{
			Status:  "failed",
			Message: "Error occured. Please try again.",
			Error:   errResponse,
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.BlockUserResponse{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})
		return
	}
	fmt.Println("********1")
	//block user
	err := h.AdminUseCase.BlockUser(&req)
	if err != nil {
		errResponse := fmt.Sprint("error while blocking user. Error: ", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.BlockUserResponse{
			Status:  "failed",
			Message: "Some error occured while blocking user.",
			Error:   errResponse,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.BlockUserResponse{
			Status:  "success",
			Message: "User blocked successfully",
			Error:   "",
		})
	}

}

// to unblock a user
// @Summary Unblock user
// @Description Unblock user
// @Tags admin
// @Accept json
// @Produce json
// @Param unblockUserRequest body requestModels.UnblockUserReq true "user"
// @Success 200 {object} response.UnblockUserResponse
// @Failure 400 {object} response.UnblockUserResponse
// @Router /admin/unblockuser [post]
func (h *AdminHandler) UnblockUser(c *gin.Context) {
	fmt.Println("Handler ::: \"Unblock user\" handler")

	//get user info from request
	var req requestModels.UnblockUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse := "error binding the request. Try again. Error:" + err.Error()
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.UnblockUserResponse{
			Status:  "failed",
			Message: "Error occured. Please try again.",
			Error:   errResponse,
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.UnblockUserResponse{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})
		return
	}

	//unblock user
	err := h.AdminUseCase.UnblockUser(&req)
	if err != nil {
		errResponse := fmt.Sprint("error while unblocking user. Error: ", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.UnblockUserResponse{
			Status:  "failed",
			Message: "Some error occured while unblocking user.",
			Error:   errResponse,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.UnblockUserResponse{
			Status:  "success",
			Message: "User unblocked successfully",
			Error:   "",
		})
	}

}

// to block a seller
// @Summary Block seller
// @Description Block seller
// @Tags admin
// @Accept json
// @Produce json
// @Param blockSellerRequest body requestModels.BlockSellerReq true "user"
// @Success 200 {object} response.BlockSellerResponse
// @Failure 400 {object} response.BlockSellerResponse
// @Router /admin/blockseller [post]
func (h *AdminHandler) BlockSeller(c *gin.Context) {
	fmt.Println("Handler ::: \"Block seller\" handler")

	//get seller info from request
	var req requestModels.BlockSellerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse := "error binding the request. Try again. Error:" + err.Error()
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.BlockSellerResponse{
			Status:  "failed",
			Message: "Error occured. Please try again.",
			Error:   errResponse,
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.BlockSellerResponse{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})
		return
	}

	//block seller
	err := h.AdminUseCase.BlockSeller(&req)
	if err != nil {
		errResponse := fmt.Sprint("error while blocking seller. Error: ", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.BlockSellerResponse{
			Status:  "failed",
			Message: "Some error occured while blocking seller.",
			Error:   errResponse,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.BlockSellerResponse{
			Status:  "success",
			Message: "Seller blocked successfully",
			Error:   "",
		})
	}

}

// to unblock a seller
// @Summary Unblock seller
// @Description Unblock seller
// @Tags admin
// @Accept json
// @Produce json
// @Param unblockSellerRequest body requestModels.UnblockSellerReq true "user"
// @Success 200 {object} response.UnblockSellerResponse
// @Failure 400 {object} response.UnblockSellerResponse
// @Router /admin/unblockseller [post]
func (h *AdminHandler) UnblockSeller(c *gin.Context) {
	fmt.Println("Handler ::: \"Unblock seller\" handler")

	//get seller info from request
	var req requestModels.UnblockSellerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse := "error binding the request. Try again. Error:" + err.Error()
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.UnblockSellerResponse{
			Status:  "failed",
			Message: "Error occured. Please try again.",
			Error:   errResponse,
		})
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Sprint("error validating the request. Try again. Error:", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.UnblockSellerResponse{
			Status:  "failed",
			Message: "#",
			Error:   errResponse,
		})
		return
	}

	//unblock seller
	err := h.AdminUseCase.UnblockSeller(&req)
	if err != nil {
		errResponse := fmt.Sprint("error while unblocking seller. Error: ", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.UnblockSellerResponse{
			Status:  "failed",
			Message: "Some error occured while unblocking seller.",
			Error:   errResponse,
		})
		return
	} else {
		c.JSON(http.StatusOK, response.UnblockSellerResponse{
			Status:  "success",
			Message: "Seller unblocked successfully",
			Error:   "",
		})
	}

}
