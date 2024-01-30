package middleware

import (
	response "MyShoo/internal/models/responseModels"
	jwttoken "MyShoo/pkg/jwt_tokens"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func SellerAuth(c *gin.Context) {
	// fmt.Println("#Middleware: SellerAuth: Entered")
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString) //
	secretKey := os.Getenv("SECRET_KEY")

	IsTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, secretKey)
	if !IsTokenValid {
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

	//checking if role is Seller
	if claims.Role != "seller" {
		fmt.Println("role is not Seller")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	seller := claims.Model

	c.Set("SellerModel", seller)
	c.Next()
}
