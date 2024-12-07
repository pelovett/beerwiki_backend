package user

import (
	"fmt"
	"net/http"

	"github.com/pelovett/beerwiki_backend/db_wrapper"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func ConfirmUser(c *gin.Context) {
	type Code struct {
		Code string `json:"code"`
	}

	var code Code

	// Parse the JSON body to get the code
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Check if the code exists in the database
	account_id, exists, err := isStringInDb(code.Code, "confirm_code")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	if exists {
		// Mark the user as confirmed in the database
		// Respond with success
		updateVerifiedAt(account_id)
		c.JSON(http.StatusOK, gin.H{"message": "Email confirmed successfully!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid confirmation code"})
	}
}

// Define a function to check if the code exists
func isStringInDb(inputString string, columnName string) (int, bool, error) {

	var accountID int
	query := fmt.Sprintf("SELECT account_id FROM users WHERE %s = $1", columnName)
	rows, err := db_wrapper.Query(query, inputString)
	if err != nil {
		return 0, false, err
	}
	defer rows.Close()

	if rows.Next() {
		// Scan the result into accountID
		err := rows.Scan(&accountID)
		if err != nil {
			return 0, false, err
		}
		// String found
		return accountID, true, nil
	}

	// No rows found
	return 0, false, nil
}

func updateVerifiedAt(userID int) error {

	_, err := db_wrapper.Query(`UPDATE users SET verified_at = NOW() WHERE account_id = $1`, userID)
	if err != nil {
		return err
	}
	return err
}
