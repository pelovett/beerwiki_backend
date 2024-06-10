package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func VerifyUser(c *gin.Context) {
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db_host := os.Getenv("DB_HOST")
	db_pass := os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=postgres "+
		"password=%s dbname=postgres sslmode=disable",
		db_host, db_pass)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	// Query the database for the user by username
	row := db.QueryRow("SELECT email, password FROM users WHERE email = $1", user.Email)

	var dbEmail, dbPassword string

	err = row.Scan(&dbEmail, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Passwords match, user is verified
	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}
