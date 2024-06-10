package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func VerifyUser(c *gin.Context) {

	cookieValue := c.GetHeader("Cookie")

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.Parse(cookieValue, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if token != nil {
		if token.Valid {
			c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("That's not even a token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				fmt.Println("Timing is everything")
			} else {
				fmt.Println("Couldn't handletoken contains an invalid number of segments this token:", err)
			}
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Token is nil:", err)
	}
}
