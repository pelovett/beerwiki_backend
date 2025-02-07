package image

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

type CompletedUploadInfo struct {
	ID string `json:"id"`
}

// Get a cloudflare creator upload URL and store temp record
func GetImageUploadURL(c *gin.Context) {
	// Check that we have user account id in state
	// accountId := c.GetInt("account_id")
	// if accountId == 0 {
	// 	c.JSON(http.StatusUnauthorized, "Missing credentials")
	// 	return
	// }
	accountId := 1

	// Check if image name was provided
	name, namePresent := c.Request.URL.Query()["name"]
	if !namePresent || len(name) < 1 || len(name[0]) < 1 || len(name[0]) > 64 {
		c.JSON(http.StatusBadRequest, "Must provide 'name' parameter")
		return
	}
	imageName := strings.ToLower(name[0])

	// Check if provided image name is unique
	result, err := db_wrapper.Query(
		"SELECT COUNT(*) != 0 AS NOT_UNIQUE FROM image_metadata WHERE image_name = $1;",
		imageName,
	)
	if err != nil || !result.Next() {
		c.JSON(http.StatusInternalServerError, "Failed to check name uniqueness")
		log.Printf("Failed to check image name: %s", err)
		return
	}
	var not_unique bool
	result.Scan(&not_unique)
	if not_unique {
		c.JSON(http.StatusBadRequest, "Name must be unique")
		return
	}

	// Get preauth upload link from cloudflare
	cfApi, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to contact image upload service")
		log.Printf("Failed to create cloudflare client: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	requireSignedUrls := false
	expiryTime := time.Now().Add(time.Minute * 30)
	uploadUrlInfo, err := cfApi.CreateImageDirectUploadURL(ctx,
		cloudflare.AccountIdentifier(os.Getenv("CLOUDFLARE_ACCOUNT_ID")),
		cloudflare.CreateImageDirectUploadURLParams{
			Version:           cloudflare.ImagesAPIVersionV2,
			RequireSignedURLs: &requireSignedUrls,
			Expiry:            &expiryTime,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to create image upload URL")
		log.Printf("Failed to create request image upload URL: %s", err)
		return
	}

	// Add id to db as in-progress upload
	_, err = db_wrapper.Exec(
		"INSERT INTO image_metadata (image_id, account_id, image_name) VALUES ($1, $2, $3);",
		uploadUrlInfo.ID,
		accountId,
		imageName,
	)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"image_metadata_image_name_key\"" {
			c.JSON(http.StatusBadRequest, "Name must be unique")
			log.Printf("Failed to save image id in postgres: |%s|", err)
		} else {
			c.JSON(http.StatusInternalServerError, "Failed to save image id in db")
			log.Printf("Failed to save image id in postgres: |%s|", err)
		}
		return
	}

	// Send the url and id back to the user
	c.JSON(http.StatusCreated, gin.H{"url": uploadUrlInfo.UploadURL, "id": uploadUrlInfo.ID})
}

// Confirm that an image has successfully been uploaded to cloudflare
func PostImageUploadComplete(c *gin.Context) {
	var uploadInfo CompletedUploadInfo
	if err := c.BindJSON(&uploadInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse request")
		log.Printf("Failed to parse image upload complete: %s", err)
		return
	}

	// Check cloudflare image status
	cfApi, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to contact image upload service")
		log.Printf("Failed to create cloudflare client: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err = cfApi.GetImage(
		ctx,
		cloudflare.AccountIdentifier(os.Getenv("CLOUDFLARE_ACCOUNT_ID")),
		uploadInfo.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to retrieve image info")
		log.Printf("Failed to get cf image info: %s", err)
		return
	}

	// Update image status in db
	result, err := db_wrapper.Exec(
		"UPDATE image_metadata SET upload_complete = TRUE WHERE image_id = $1",
		uploadInfo.ID,
	)
	numRows, _ := result.RowsAffected()
	if err != nil || numRows == 0 {
		c.JSON(http.StatusInternalServerError, "Failed to update image status")
		if numRows == 0 {
			log.Printf("Couldn't find id: %s", uploadInfo.ID)
		} else {
			log.Printf("Error: %s", err)
		}
		return
	}
}

// Confirm that an image has successfully been uploaded to cloudflare
func GetImageURL(c *gin.Context) {
	// Get image name
	imageName := strings.ToLower(c.Param("name"))
	if imageName == "" {
		c.JSON(http.StatusBadRequest, "No image name provided")
		return
	}

	// Check if provided image name is unique
	results, err := db_wrapper.Query(
		"SELECT image_id FROM image_metadata WHERE image_name = $1 limit 1;",
		imageName,
	)
	if err != nil || !results.Next() {
		c.JSON(http.StatusNotFound, "Failed to retrieve image info")
		log.Printf("Failed to get image name: %s", err)
		return
	}
	var imageId string
	results.Scan(&imageId)

	c.Writer.Header().Set("Cache-Control", "public, max-age=600, immutable")
	c.JSON(
		http.StatusOK,
		gin.H{"url": fmt.Sprintf(
			"https://imagedelivery.net/%s/%s/public",
			os.Getenv("CLOUDFLARE_ACCOUNT_HASH"),
			imageId,
		),
		},
	)
}
