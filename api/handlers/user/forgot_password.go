package user

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

func ForgotPassword(c *gin.Context) {
	// Define the expected request payload structure
	type ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"` // Updated to `json:"email"` for JSON payload
	}

	var forgotPasswordRequest ForgotPasswordRequest

	// Parse the JSON body to get the email
	if err := c.ShouldBindJSON(&forgotPasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Check if the email exists in the database
	accountId, exists, err := isStringInDb(forgotPasswordRequest.Email, "email")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	if exists {
		resetCode := uuid.NewString()
		err := setResetCode(resetCode, accountId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		} else {
			resetURL, err := generateURL(resetCode, "/user/change-password/")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Generating URL failed"})
			} else {
				err := sendPasswordResetEmail(forgotPasswordRequest.Email, resetURL)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Sending Email"})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "Email confirmed successfully!"})
				}
			}
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}
}

func setResetCode(resetCode string, accountId int) error {
	_, err := db_wrapper.Exec(
		"UPDATE users SET reset_code = $1, reset_at = now() AT TIME ZONE 'UTC' WHERE account_id = $2;",
		resetCode, accountId,
	)
	return err

}

func sendPasswordResetEmail(userEmail string, resetLink string) error {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return err
	}

	svc := ses.New(sess)

	// Create HTML email body
	htmlBody := fmt.Sprintf(`
        <html>
            <body>
                <h1 style="color: #4CAF50;">HopWiki!</h1>
                <p>Please click the link below to reset your password:</p>
				<p>This link will expire in 24 hours</p>
                <p>
                    <a href="%s" style="padding: 10px 15px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Reset Password</a>
                </p>
                <p>If you did request a password reset, please ignore this email.</p>
                <p>Best regards,<br>The HopWiki Team</p>
            </body>
        </html>
    `, resetLink)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(userEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(fmt.Sprintf("Reset your password by clicking the link: %s", resetLink)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Password Reset"),
			},
		},
		Source: aws.String("account-services@hopwiki.org"),
	}

	_, err = svc.SendEmail(input)
	return err
}
