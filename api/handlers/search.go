package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

func SearchBeerByName(c *gin.Context) {
	params := c.Request.URL.Query()
	query, queryPresent := params["q"]
	if !queryPresent || len(query) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	rows, err := db_wrapper.Query(
		"SELECT name, url_name, similarity(name, $1) as sim FROM beer ORDER BY similarity(name, $1) DESC LIMIT 5;",
		query[0],
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Failed to search db: %s", err)
		return
	}

	// Iterate over the rows
	results := []map[string]string{}
	for rows.Next() {
		var name string
		var urlName string
		var similarity float32

		// Scan the values from the current row into variables
		if err := rows.Scan(&name, &urlName, &similarity); err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("Failed to parse search results: %s", err)
			return
		}

		// Skip results that dont really match
		if similarity <= 0.0001 {
			continue
		}

		results = append(results, map[string]string{
			"name":             name,
			"url_name":         urlName,
			"similarity_score": fmt.Sprintf("%f", similarity),
		})
	}

	// If no results are found, then say so
	if len(results) > 0 {
		c.JSON(http.StatusOK, map[string]any{
			"results": results,
		})
	} else {
		c.Status(http.StatusNotFound)
	}
}
