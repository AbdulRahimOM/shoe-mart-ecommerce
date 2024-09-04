package jwttoken

import (
	"fmt"
	"time"

	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/config"

	jwt "github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Model interface{}
	Role  string
	jwt.RegisteredClaims
}

func GenerateToken(role string, myClaims interface{}, validityDuration time.Duration) (string, error) {

	//create a custom claim
	claims := &CustomClaims{
		myClaims,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validityDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.SecretKey))

	// fmt.Println(tokenString, err)
	return tokenString, err
}
func IsTokenValid(tokenString string, secretKey string) (bool, interface{}) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println("token is invalid or error in parsing")
		return false, nil
	}
	if !token.Valid {
		fmt.Println("token is invalid")
		return false, nil
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		//check if token expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			fmt.Println("token expired")
			return false, nil
		}

		return true, claims
	} else {
		fmt.Println("Error while decoding token")
		return false, nil
	}
}
