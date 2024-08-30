package user

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

type JwtClaims struct {
	Email     string `json:"email"`
	AccountId int    `json:"account_id"`
	jwt.StandardClaims
}

// Check if cookie is valid and return account id if so
// Should maybe check if user exists in db but whatever
func CheckValidCookie(cookie string) int {

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if token == nil {
		fmt.Println("Token is nil: ", err)
		return 0
	}

	// If token has expired, then its invalid
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return 0
	}

	if token.Valid {

		return claims.AccountId
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle token contains an invalid number of segments this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	return 0
}

func VerifyUser(c *gin.Context) {
	// If account isn't in state, we'll get 0 back
	var verified bool
	accountId := c.GetInt("account_id")
	if accountId != 0 {
		rows, err := db_wrapper.Query(`SELECT verified_at IS NOT NULL AS is_valid FROM users WHERE account_id = $1`, accountId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth cookie provided"})
		}
		if rows.Next() {
			// Scan the result into accountID
			err := rows.Scan(&verified)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not verified"})
			} else if verified {
				c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
			}

		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth cookie provided"})
	}
}
