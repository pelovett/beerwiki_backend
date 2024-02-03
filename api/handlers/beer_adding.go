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
	fmt.Printf("The id is %s\n", id)
	result, err := db.Query("SELECT name FROM beer WHERE id=$1 limit 1;", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Beer not found",
		})
		log.Println(`Failed to run query: ${err}`)
		return
	}
	var beer string
	if result.Next() {
		if err := result.Scan(&beer); err != nil {
			log.Println(`Failed to parse beer string id: %s`, id)
		}
	}

	if beer != "" {
		fmt.Printf("Returning beer: %s\n", beer)
		c.JSON(http.StatusOK, gin.H{
			"message": beer,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Beer not found",
		})
	}

}
