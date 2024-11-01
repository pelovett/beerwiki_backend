package user

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
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

	err = createUserAccount(createUserRequest.Email, createUserRequest.Username, hashedPassword, createUserRequest.CreatedAt, confirmCode)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account created successfully for user: %s", createUserRequest.Username)})
		url, err := generateConfirmationURL(confirmCode)
		if err != nil {
			log.Println("Error creating user account:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {

			sendConfirmEmail(createUserRequest.Email, url)
			if err != nil {
				log.Println("Error sending email:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func createUserAccount(email string, user_name string, password string, createdAt time.Time, confirm_code string) error {
	_, err := db_wrapper.Exec(
		"INSERT INTO users (email, user_name, password, created_at, confirm_code) VALUES ($1, $2, $3, now() AT TIME ZONE 'UTC', $4)",
		email,
		user_name,
		password,
		confirm_code,
	)

	if err != nil {
		// Use a switch to handle different types of unique constraint violations
		log.Println("Error creating user account:", err)
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // Unique violation error code for PostgreSQL
				// Further check if it's the unique constraint on `email` or `user_name`
				if pqErr.Constraint == "users_email_key" {
					return fmt.Errorf("email already exists")
				}
				if pqErr.Constraint == "users_user_name_key" {
					return fmt.Errorf("user name already exists")
				}
			case "23502": // Not-null constraint violation
				if pqErr.Column == "email" {
					return fmt.Errorf("email cannot be null")
				}
				if pqErr.Column == "user_name" {
					return fmt.Errorf("user name cannot be null")
				}
				if pqErr.Column == "password" {
					return fmt.Errorf("password cannot be null")
				}
				if pqErr.Column == "created_at" {
					return fmt.Errorf("created_at cannot be null")
				}
			case "22001": // Value too long for type
				if pqErr.Column == "confirm_code" {
					return fmt.Errorf("confirm code is too long (max 255 characters)")
				}
			default:
				// Handle other types of database errors
				return fmt.Errorf("database error: %s", pqErr.Message)
			}
		}
		// Return the original error if it's not related to a constraint violation
		return err
	}

	return nil
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
func generateConfirmationURL(confirmCode string) (string, error) {
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
func sendConfirmEmail(userEmail, confirmationLink string) error {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return err
	}

	svc := ses.New(sess)

	// Create HTML email body
	htmlBody := fmt.Sprintf(`
        <html>
            <body>
                <h1 style="color: #4CAF50;">Welcome to HopWiki!</h1>
                <p>Thank you for signing up. Please confirm your account by clicking the link below:</p>
                <p>
                    <a href="%s" style="padding: 10px 15px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Confirm My Account</a>
                </p>
                <p>If you did not sign up, please ignore this email.</p>
                <p>Best regards,<br>The HopWiki Team</p>
            </body>
        </html>
    `, confirmationLink)

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
					Data:    aws.String(fmt.Sprintf("Please confirm your account by clicking the link: %s", confirmationLink)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Account Confirmation"),
			},
		},
		Source: aws.String("account-services@hopwiki.org"),
	}

	_, err = svc.SendEmail(input)
	return err
}
