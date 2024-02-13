package ordermanagementHandlers

import (
	"MyShoo/internal/domain/entities"
	requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishListHandler struct {
	wishListUseCase usecaseInterface.IWishListsUC
}

func NewWishListHandler(wishListUseCase usecaseInterface.IWishListsUC) *WishListHandler {
	return &WishListHandler{wishListUseCase: wishListUseCase}
}

// create new wishlist
func (h *WishListHandler) CreateWishList(c *gin.Context) {
	var req *requestModels.CreateWishListReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error: %v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
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
func (h *WishListHandler) AddToWishList(c *gin.Context) {

	var req *requestModels.AddToWishListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error: %v", err)
		fmt.Println(errResponse)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
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
func (h *WishListHandler) RemoveFromWishList(c *gin.Context) {

	var req *requestModels.RemoveFromWishListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedSME("Error binding request. Try Again", err))
		return
	}

	//validate request
	if err := requestValidation.ValidateRequest(req); err != nil {
		errResponse := fmt.Errorf("error validating the request. Try again. Error:%v", err)
		c.JSON(http.StatusBadRequest, response.FailedSME("Error validating request. Try Again", errResponse))
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
