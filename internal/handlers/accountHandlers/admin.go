package accounthandler

import (
	"net/http"

	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	usecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/interface"
	requestValidation "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/validation"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminUseCase usecase.IAdminUC
}

func NewAdminHandler(useCase usecase.IAdminUC) *AdminHandler {
	return &AdminHandler{AdminUseCase: useCase}
}

// to get admin login page
// @Summary Get admin login page
// @Description Get admin login page
// @Tags Admin/Session
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} string
// @Router /admin/login [get]
func (h *AdminHandler) GetAdminLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

// to login admin
// @Summary Login admin
// @Description Login admin
// @Tags Admin/Session
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param adminSignInReq body req.AdminSignInReq{} true "Admin Sign In Request"
// @Success 200 {object} response.SMT{}
// @Failure 400 {object} response.SME{}
// @Router /admin/login [post]
func (h *AdminHandler) PostLogIn(c *gin.Context) {

	var req request.AdminSignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	token, err := h.AdminUseCase.SignIn(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.SMT{
			Token: *token,
		})
	}
}

// to get users list
// @Summary Get users list
// @Description Get users list
// @Tags Admin/Account_Management/Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetUsersListResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/userslist [get]
func (h *AdminHandler) GetUsersList(c *gin.Context) {

	usersList, err := h.AdminUseCase.GetUsersList()
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.GetUsersListResponse{
			UsersList: *usersList,
		})
	}
}

// to get sellers list
// @Summary Get sellers list
// @Description Get sellers list
// @Tags Admin/Account_Management/Sellers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetSellersListResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/sellerslist [get]
func (h *AdminHandler) GetSellersList(c *gin.Context) {

	sellersList, err := h.AdminUseCase.GetSellersList()
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.GetSellersListResponse{
			SellersList: *sellersList,
		})
	}
}

// to block a user
// @Summary Block user
// @Description Block user
// @Tags Admin/Account_Management/Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param unblockUserRequest body req.BlockUserReq{} true "user"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/blockuser [post]
func (h *AdminHandler) BlockUser(c *gin.Context) {

	//get user info from request
	var req request.BlockUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//block user
	err := h.AdminUseCase.BlockUser(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("User blocked"))
	}

}

// to unblock a user
// @Summary Unblock user
// @Description Unblock user
// @Tags Admin/Account_Management/Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param unblockUserRequest body req.UnblockUserReq true "user"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/unblockuser [post]
func (h *AdminHandler) UnblockUser(c *gin.Context) {

	//get user info from request
	var req request.UnblockUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//unblock user
	err := h.AdminUseCase.UnblockUser(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("User unblocked"))
	}

}

// to block a seller
// @Summary Block seller
// @Description Block seller
// @Tags Admin/Account_Management/Sellers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param blockSellerRequest body req.BlockSellerReq true "user"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/blockseller [post]
func (h *AdminHandler) BlockSeller(c *gin.Context) {

	//get seller info from request
	var req request.BlockSellerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//block seller
	err := h.AdminUseCase.BlockSeller(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("Seller blocked"))
	}

}

// to unblock a seller
// @Summary Unblock seller
// @Description Unblock seller
// @Tags Admin/Account_Management/Sellers
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param unblockSellerRequest body req.UnblockSellerReq true "user"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /admin/unblockseller [post]
func (h *AdminHandler) UnblockSeller(c *gin.Context) {

	//get seller info from request
	var req request.UnblockSellerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	//unblock seller
	err := h.AdminUseCase.UnblockSeller(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessSM("Seller unblocked"))
	}

}

// VerifySeller
func (h *AdminHandler) VerifySeller(c *gin.Context) {
	var req request.VerifySellerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnBindingReq(err))
		return
	}

	//validation
	if err := requestValidation.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrOnFormValidation(&err))
		return
	}

	err := h.AdminUseCase.VerifySeller(&req)
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Seller verified"))
}

// ReloadConfig
// @Summary Reload config
// @Description Reload config
// @Tags Admin/System_Related/Config
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.SM{} @Example {"status": "success", "message": "Config reloaded"}
// @Failure 400 {object} response.SME{}
// @Router /admin/system/restart-Configuration [get]
func (h *AdminHandler) RestartConfig(c *gin.Context) {
	err := h.AdminUseCase.RestartConfig()
	if err != nil {
		c.JSON(err.StatusCode, response.FromError(err))
		return
	}
	c.JSON(http.StatusOK, response.SuccessSM("Config reloaded"))
}
