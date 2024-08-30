package user

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
	// GenerateUUID generates a UUID which can be used as a unique confirmation code
)

func CreateUser(c *gin.Context) {
	type CreateUserRequest struct {
		Email     string    `form:"email" binding:"required,email"`
		Username  string    `form:"username" binding:"required"`
		Password  string    `form:"password" binding:"required"`
		CreatedAt time.Time `form:"created_at" time_format:"2006-01-02T15:04:05Z07:00"`
	}

	var createUserRequest CreateUserRequest

	if err := c.ShouldBind(&createUserRequest); err != nil {
		log.Println("Binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(createUserRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	confirmCode := uuid.NewString()

	if createUserAccount(createUserRequest.Email, createUserRequest.Username, hashedPassword, createUserRequest.CreatedAt, confirmCode) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account created successfully for user: %s", createUserRequest.Username)})
		url, err := GenerateConfirmationURL(confirmCode)
		if err != nil {
			log.Println("Error creating user account:", err)
		}
		print(url)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
	}
}

func createUserAccount(email string, user_name string, password string, createdAt time.Time, confirm_code string) bool {
	output, err := db_wrapper.Exec(
		"INSERT INTO users (email, user_name, password, created_at, confirm_code) VALUES ($1, $2, $3, now() AT TIME ZONE 'UTC', $4)",
		email,
		user_name,
		password,
		confirm_code,
	)
	if err != nil {
		log.Println("Error creating user account:", err)
		return false
	}

	log.Println(output)

	return true
}

func hashPassword(password string) (string, error) {
	// Generate a salt with a cost of 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GenerateConfirmationURL creates a URL for email confirmation with the provided confirmation code
func GenerateConfirmationURL(confirmCode string) (string, error) {
	// Get the base URL from an environment variable
	baseURL := os.Getenv("SERVER_ADDRESS")
	if baseURL == "" {
		return "", fmt.Errorf("SERVER_ADDRESS environment variable is not set")
	}

	confirmURL := baseURL + "/user/confirmation/"

	// Parse the base URL
	u, err := url.Parse(confirmURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %v", err)
	}

	// Add the query parameter for the confirmation code
	q := u.Query()
	q.Set("code", confirmCode)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
