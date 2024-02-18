package middleware

import (
	"MyShoo/internal/config"
	response "MyShoo/internal/models/responseModels"
	jwttoken "MyShoo/pkg/jwt"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString) //

	IsTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, config.SecretKey)
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
	c.Next()
}
