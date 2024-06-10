package beer

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

// Gets beer name from database if it exists.
func GetBeer(c *gin.Context) {
	id := c.Param("id")
	result, err := db_wrapper.Query(
		"SELECT id, name, url_name, page_ipa_ml FROM beer WHERE id=$1 limit 1;",
		id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Beer not found",
		})
		log.Printf("Failed to run query: %s", err)
		return
	}

	var foundBeer beer
	if result.Next() {
		if err := result.Scan(&foundBeer.ID, &foundBeer.Name, &foundBeer.URLName, &foundBeer.PageIPAML); err != nil {
			log.Printf("Failed to parse beer string id: %s", id)
		}
	}

	if foundBeer.Name != "" {
		c.JSON(http.StatusOK, foundBeer)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Beer not found",
		})
	}

}

// Gets beer name from database if it exists.
func GetBeerByUrlName(c *gin.Context) {

	url_name := c.Param("name")
	result, err := db_wrapper.Query(
		"SELECT id, name, url_name, page_ipa_ml FROM beer WHERE url_name=$1 limit 1;",
		url_name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Beer not found",
		})
		log.Printf("Failed to run query: %s", err)
		return
	}

	var foundBeer beer
	if result.Next() {
		if err := result.Scan(&foundBeer.ID, &foundBeer.Name, &foundBeer.URLName, &foundBeer.PageIPAML); err != nil {
			log.Printf("Failed to parse beer string url_name: %s", url_name)
		}
	}

	if foundBeer.Name != "" {
		c.JSON(http.StatusOK, foundBeer)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Beer not found",
		})
	}

}
