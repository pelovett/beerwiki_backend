package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

func CreateUser(c *gin.Context) {
	var formData struct {
		Email     string    `from:"email" binding:"required,email"`
		Username  string    `form:"user_name" binding:"required"`
		Password  string    `form:"password" binding:"required"`
		CreatedAt time.Time `form:"created_at" time_format:"2006-01-02T15:04:05Z07:00"`
	}

	if err := c.ShouldBind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hashedPassword, err := hashPassword(formData.Password)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
	// 	return
	// }

	if createUserAccount(formData.Email, formData.Username, formData.Password, formData.CreatedAt) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account created successfully for user: %s", formData.Username)})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
	}
}

func createUserAccount(email string, user_name string, password string, createdAt time.Time) bool {
	output, err := db_wrapper.Exec(
		"INSERT INTO users (email, user_name, password, created_at) VALUES ($1, $2, $3, now() AT TIME ZONE 'UTC')",
		email,
		user_name,
		password,
	)
	if err != nil {
		log.Println("Error creating user account:", err)
		return false
	}

	log.Println(output)

	return true
}
