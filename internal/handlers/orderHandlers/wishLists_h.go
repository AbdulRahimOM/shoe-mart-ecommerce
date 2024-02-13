package orderhandler

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
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishListHandler struct {
	wishListUseCase usecase.IWishListsUC
}

func NewWishListHandler(wishListUseCase usecase.IWishListsUC) *WishListHandler {
	return &WishListHandler{wishListUseCase: wishListUseCase}
}

// create new wishlist
// @Summary Create wishlist
// @Description Create wishlist
// @Tags User/WishList
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param createWishListReq body req.CreateWishListReq{} true "Create WishList Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /createwishlist [post]
func (h *WishListHandler) CreateWishList(c *gin.Context) {

	var req *request.CreateWishListReq
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
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error creating wishlist. Try Again", err))
		return
	}

	//create wishlist
	err = h.wishListUseCase.CreateWishList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error creating wishlist. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Created wishlist successfully"))
}

// add to wishlist
// @Summary Add to wishlist
// @Description Add to wishlist
// @Tags User/WishList
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param addToWishListReq body req.AddToWishListReq{} true "Add to WishList Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /addtowishlist [post]
func (h *WishListHandler) AddToWishList(c *gin.Context) {

	var req *request.AddToWishListReq
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
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding to wishlist. Try Again", err))
		return
	}

	//add to wishlist
	err = h.wishListUseCase.AddToWishList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error adding to wishlist. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Added to wishlist successfully"))
}

// remove from wishlist
// @Summary Remove from wishlist
// @Description Remove from wishlist
// @Tags User/WishList
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param removeFromWishListReq body req.RemoveFromWishListReq{} true "Remove from WishList Request"
// @Success 200 {object} response.SM{}
// @Failure 400 {object} response.SME{}
// @Router /removefromwishlist [delete]
func (h *WishListHandler) RemoveFromWishList(c *gin.Context) {

	var req *request.RemoveFromWishListReq
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
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error removing from wishlist. Try Again", err))
		return
	}

	//remove from wishlist
	err = h.wishListUseCase.RemoveFromWishList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error removing from wishlist. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessSM("Removed from wishlist successfully"))
}

// GetAllWishLists
// @Summary Get all wishlists
// @Description Get all wishlists
// @Tags User/WishList
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.GetAllWishListsResponse{}
// @Failure 400 {object} response.SME{}
// @Router /mywishlists [get]
func (h *WishListHandler) GetAllWishLists(c *gin.Context) {

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting wishlists. Try Again", err))
		return
	}

	//get all wishlists
	var wishLists *[]entities.WishList
	wishLists, totalCount, err := h.wishListUseCase.GetAllWishLists(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting wishlists. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.GetAllWishListsResponse{
		Status:  "success",
		Message: "Got wishlists successfully",

		WishLists:  *wishLists,
		TotalCount: totalCount,
	})
}

// GetWishListByID
// @Summary Get wishlist by id
// @Description Get wishlist by id
// @Tags User/WishList
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query string true "WishList ID"
// @Success 200 {object} response.GetWishListByIDResponse{}
// @Failure 400 {object} response.SME{}
// @Router /wishlist [get]
func (h *WishListHandler) GetWishListByID(c *gin.Context) {

	//get query params
	wishListID := c.Query("id")

	//validate params
	if wishListID == "" {
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Wishlist id is empty",
			Error:   "wishlist id is empty",
		})
		return
	}

	//convert string to uint
	wishListIDUint, err := strconv.ParseUint(wishListID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error converting wishlist id to uint. Try Again", err))
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting wishlist. Try Again", err))
		return
	}

	//get wishlist by id
	wishListName, wishListItems, totalCount, err := h.wishListUseCase.GetWishListByID(userID, uint(wishListIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedSME("Error getting wishlist. Try Again", err))
		return
	}

	c.JSON(http.StatusOK, response.GetWishListByIDResponse{
		Status:       "success",
		Message:      "Got wishlist successfully",
		Error:        "",
		WishListName: *wishListName,
		WishItems:    *wishListItems,
		TotalCount:   totalCount,
	})
}
