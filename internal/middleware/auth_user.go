package middleware

import (
	response "MyShoo/internal/models/responseModels"
	jwttoken "MyShoo/pkg/jwt"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserAuth(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	// fmt.Println("tokenString: ", tokenString) //
	secretKey := os.Getenv("SECRET_KEY")

	isTokenValid, tokenClaims := jwttoken.IsTokenValid(tokenString, secretKey)
	if !isTokenValid {
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

	//checking if role is user
	if claims.Role != "user" {
		fmt.Println("role is not user")
		c.JSON(http.StatusUnauthorized, response.UnauthorizedAccess)
		c.Abort()
		return
	}

	user := claims.Model

	c.Set("UserModel", user)
	// fmt.Println("usermodel in c.context:", c.GetString("UserModel"))
	c.Next()
}

func NotLoggedOutCheck(c *gin.Context) {
	fmt.Println("Handler ::: check for already signed-in before going to Next()")
	//if aleady logged in
	//...redirect to home
	//...abort() if required
	//else
	fmt.Println("not logged in. So, allowing into sign-up page")
	c.Next()
}
