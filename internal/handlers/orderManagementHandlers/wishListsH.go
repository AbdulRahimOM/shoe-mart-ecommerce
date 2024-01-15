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
	fmt.Println("Handler ::: create wishlist handler")
	var req *requestModels.CreateWishListReq

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
			Message: "Error creating wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//create wishlist
	err = h.wishListUseCase.CreateWishList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error creating wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Wishlist created successfully",
		Error:   "",
	})
}

// add to wishlist
func (h *WishListHandler) AddToWishList(c *gin.Context) {
	fmt.Println("Handler ::: add to wishlist handler")
	var req *requestModels.AddToWishListReq

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
			Message: "Error adding to wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//add to wishlist
	err = h.wishListUseCase.AddToWishList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error adding to wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Added to wishlist successfully",
		Error:   "",
	})
}

// remove from wishlist
func (h *WishListHandler) RemoveFromWishList(c *gin.Context) {
	fmt.Println("Handler ::: remove from wishlist handler")
	var req *requestModels.RemoveFromWishListReq

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
			Message: "Error removing from wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//remove from wishlist
	err = h.wishListUseCase.RemoveFromWishList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error removing from wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SME{
		Status:  "success",
		Message: "Removed from wishlist successfully",
		Error:   "",
	})
}

// GetAllWishLists
func (h *WishListHandler) GetAllWishLists(c *gin.Context) {
	fmt.Println("Handler ::: get all wishlists handler")

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting wishlists. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//get all wishlists
	var wishLists *[]entities.WishList
	wishLists, totalCount, err := h.wishListUseCase.GetAllWishLists(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting wishlists. Try Again",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetAllWishListsResponse{
		Status:     "success",
		Message:    "Got wishlists successfully",
		Error:      "",
		WishLists:  *wishLists,
		TotalCount: totalCount,
	})
}

// GetWishListByID
func (h *WishListHandler) GetWishListByID(c *gin.Context) {
	fmt.Println("Handler ::: get wishlist by id handler")

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
		c.JSON(http.StatusBadRequest, response.SME{
			Status:  "failed",
			Message: "Error parsing wishlist id. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//get userID from token
	userID, err := tools.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting wishlist. Try Again",
			Error:   err.Error(),
		})
		return
	}

	//get wishlist by id
	wishListName, wishListItems, totalCount, err := h.wishListUseCase.GetWishListByID(userID, uint(wishListIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.SME{
			Status:  "failed",
			Message: "Error getting wishlist. Try Again",
			Error:   err.Error(),
		})
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
