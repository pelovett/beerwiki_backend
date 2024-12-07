package user

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// logoutUser clears the session cookie
func LogoutUser(c *gin.Context) {
	server_address := os.Getenv("SERVER_ADDRESS")
	u, _ := url.Parse(server_address)
	host_name := u.Hostname()

	_, err := c.Cookie("login_cookie")
	if err != nil {
		// Cookie is not present, return 401 or appropriate error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}

	// If cookie is present, proceed to remove it
	c.SetCookie(
		"login_cookie", // Cookie name
		"",             // Clear the value
		-1,             // Set MaxAge to -1 to expire the cookie
		"/",            // Path
		host_name,      // Domain
		false,          // Secure flag
		false,          // HttpOnly flag
	)

	// Respond with a successful logout
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})

}
