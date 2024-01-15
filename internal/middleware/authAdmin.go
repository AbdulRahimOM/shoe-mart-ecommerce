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

func AdminAuth(c *gin.Context) {
	// fmt.Println("#Middleware: AdminAuth: Entered")
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

	//checking if role is Admin
	if claims.Role != "admin" {
		fmt.Println("role is not admin")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	admin := claims.Model

	c.Set("AdminModel", admin)
	fmt.Println("Admin model in c.context:", c.GetString("AdminModel"))
	fmt.Println("#Middleware Exited. Result: Admin Authenticated")
	c.Next()
}
