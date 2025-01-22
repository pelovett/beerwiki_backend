package image

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

func GetImage(c *gin.Context) {
	image_name := c.Param("image")
	result, err := db_wrapper.Query(
		"SELECT image_id, image_name FROM image_metadata WHERE image_name = $1 limit 1;",
		image_name,
	)

	if err != nil || !result.Next() {
		c.JSON(http.StatusNotFound, "Failed to retrieve image info")
		log.Printf("Failed to retrieve image: %s with error %s", image_name, err)
	}

}
