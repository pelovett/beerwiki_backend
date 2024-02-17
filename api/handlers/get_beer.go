package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Gets beer name from database if it exists.
func GetBeer(c *gin.Context) {
	db_host := os.Getenv("DB_HOST")
	db_pass := os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=postgres "+
		"password=%s dbname=postgres sslmode=disable",
		db_host, db_pass)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	id := c.Param("id")
	result, err := db.Query("SELECT id, name, url_name, page_ipa_ml FROM beer WHERE id=$1 limit 1;", id)

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
	db_host := os.Getenv("DB_HOST")
	db_pass := os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=postgres "+
		"password=%s dbname=postgres sslmode=disable",
		db_host, db_pass)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	url_name := c.Param("name")
	result, err := db.Query(
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
