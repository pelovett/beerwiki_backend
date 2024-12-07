package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

func ChangePassword(c *gin.Context) {
	// Define the expected request payload structure
	type ChangePasswordRequest struct {
		Password string `json:"password"`
		Code     string `json:"code"`
	}

	var changePasswordRequest ChangePasswordRequest

	// Parse the JSON body to get the email
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	userId, exists, err := isStringInDb(changePasswordRequest.Code, "reset_code")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	resetValid := false
	if exists {
		resetValid = checkTimeOfCol(userId, "reset_at")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password reset link"})
		return
	}

	hashedPassword, err := hashPassword(changePasswordRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	if resetValid {
		_, err = db_wrapper.Query(`UPDATE users SET password = $1, reset_at = NULL WHERE account_id = $2`, hashedPassword, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset link expired"})

	}

}

func checkTimeOfCol(userId int, colName string) bool {
	var timeCol time.Time

	// Use db_wrapper's method to run the query and get results (assuming db_wrapper returns *sql.Rows)
	rows, err := db_wrapper.Query(`SELECT `+colName+` FROM users WHERE account_id = $1`, userId)
	if err != nil {
		return false
	}
	defer rows.Close() // Make sure to close the rows

	// Check if a row is returned
	if rows.Next() {
		// Scan the value into timeCol
		err := rows.Scan(&timeCol)
		if err != nil {
			return false
		}
	} else {
		// No rows returned
		return false
	}

	// Get the current time
	currentTime := time.Now()

	// Calculate the time difference
	timeDifference := currentTime.Sub(timeCol)
	return timeDifference < 24*time.Hour
}
