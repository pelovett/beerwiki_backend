package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelovett/beerwiki_backend/api/handlers/user"
)

// Validate cookie and set account id state
func Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Try grabbing cookie from request
		cookie, err := c.Cookie("login_cookie")
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusUnauthorized)
			return
		}

		// Check if cookie is valid
		potentialAccountId := user.CheckValidCookie(cookie)
		if potentialAccountId == 0 {
			c.Status(http.StatusUnauthorized)
			return
		}

		// Store account id in state for downstream handlers
		c.Set("account_id", potentialAccountId)
		c.Next()
	}
}
