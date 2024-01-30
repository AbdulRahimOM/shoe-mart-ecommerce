package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func UserAwaitingVerification(c *gin.Context) {
	userModel, ok := c.Get("UserModel")
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "UserModel not found"})
		return
	}
	userStatus := userModel.(map[string]interface{})["status"].(string)
	fmt.Println("userStatus=", userStatus)
	if userStatus != "not verified" {
		c.AbortWithStatusJSON(403, gin.H{"error": "User is already verified. No need to verify again"})
		return
	}

	c.Next()
}
func SellerAwaitingVerification(c *gin.Context) {
	sellerModel, ok := c.Get("SellerModel")
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "SellerModel not found",
		})
		return
	}
	sellerStatus := sellerModel.(map[string]interface{})["status"].(string)
	fmt.Println("sellerStatus=", sellerStatus)
	if sellerStatus != "not verified" {
		c.AbortWithStatusJSON(403, gin.H{
			"error": "Seller is already verified. No need to verify again",
		})
		return
	}
	c.Next()
}
