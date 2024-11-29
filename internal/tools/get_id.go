package tools

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, error) {
	userModel, ok := c.Get("UserModel")
	if !ok {
		fmt.Println("UserModel not found in context")
		return 0, errors.New("UserModel not found in context")
	}
	userID, ok := userModel.(map[string]interface{})["id"].(float64)
	if !ok {
		fmt.Println("UserModel: ", userModel)
		return 0, errors.New("User ID not found in context")
	}
	return uint(userID), nil
}

func GetSellerID(c *gin.Context) (uint, error) {
	sellerModel, ok := c.Get("SellerModel")
	if !ok {
		fmt.Println("SellerModel not found in context")
		return 0, errors.New("SellerModel not found in context")
	}
	sellerID := uint(sellerModel.(map[string]interface{})["id"].(float64))
	return sellerID, nil
}
