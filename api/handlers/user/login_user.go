package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/pelovett/beerwiki_backend/db_wrapper"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func LoginUser(c *gin.Context) {
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the URL
	server_address := os.Getenv("SERVER_ADDRESS")
	u, err := url.Parse(server_address)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Query the database for the user by username
	rows, err := db_wrapper.Query("SELECT email, password, account_id FROM users WHERE email = $1", user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !rows.Next() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	var dbEmail, dbPassword string
	var accountId int
	err = rows.Scan(&dbEmail, &dbPassword, &accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Invalid email")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password)); err != nil {
		log.Println("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	// Passwords match, user is verified
	// Token expires in 24 hours
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      user.Email,
		"account_id": accountId,
		"exp":        time.Now().UTC().Add(24 * time.Hour).Unix(),
	})

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println(err)
		return
	}

	// Get the hostname (without port)
	host_name := u.Hostname()
	c.SetCookie("login_cookie", tokenString, 3600*24, "/", host_name, false, false)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
