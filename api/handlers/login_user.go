package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func LoginUser(c *gin.Context) {
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user User

	secretKey := []byte(os.Getenv("SECRET_KEY"))

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println(err)
		return
	}

	c.SetCookie("login_cookie", tokenString, 3600*24, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
