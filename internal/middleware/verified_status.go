package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Model interface{}
	Role  string
	jwt.RegisteredClaims
}

func VerifyUserStatus(c *gin.Context) {
	// fmt.Println("entered verify user status check")
	userModel, ok := c.Get("UserModel")
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "UserModel not found"})
		c.Abort()
		return
	}

	userStatus, ok := userModel.(map[string]interface{})["status"].(string)
	if !ok {
		fmt.Println("userModel: ", userModel)
		c.AbortWithStatusJSON(500, gin.H{"error": "User status not found"})
		c.Abort()
		return
	}

	if userStatus == "not verified" {
		c.AbortWithStatusJSON(403, gin.H{"error": "User not verified"})
		c.Abort()
		return
	} else if userStatus == "blocked" {
		c.AbortWithStatusJSON(403, gin.H{"error": "User blocked"})
		c.Abort()
		return
	}
	c.Next()
}

func VerifySellerStatus(c *gin.Context) {
	sellerModel, ok := c.Get("SellerModel")
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "SellerModel not found"})
		c.Abort()
		return
	}
	sellerStatus,ok := sellerModel.(map[string]interface{})["status"].(string)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "Seller status not found"})
		c.Abort()
		return
	}
	if sellerStatus == "not verified" {
		c.AbortWithStatusJSON(403, gin.H{"error": "Seller not verified"})
		c.Abort()
		return
	} else if sellerStatus == "blocked" {
		c.AbortWithStatusJSON(403, gin.H{"error": "Seller blocked"})
		c.Abort()
		return
	}
	c.Next()
}
