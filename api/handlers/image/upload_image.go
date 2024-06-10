package image

import (
	"context"
	"log"
	"net/http"
	"os"
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

	// Add id to db as in progress upload
	_, err = db_wrapper.Exec("INSERT INTO image_metadata (image_id) VALUES ($1);", uploadUrlInfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to save image id in db")
		log.Printf("Failed to save image id in postgres: %s", err)
		return
	}

	// Send the url back to the user
	c.JSON(http.StatusAccepted, gin.H{"url": uploadUrlInfo.UploadURL})
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

	_, err = cfApi.GetImage(ctx,
		cloudflare.AccountIdentifier(os.Getenv("CLOUDFLARE_ACCOUNT_ID")),
		uploadInfo.ID)
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
