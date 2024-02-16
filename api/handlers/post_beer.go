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

type beer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func PostBeer(c *gin.Context) {

	var newBeer beer

	if err := c.BindJSON(&newBeer); err != nil {
		log.Printf("Failed to bind inputted json to beer with err %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println(newBeer)

	if err := addBeerDB(&newBeer); err != nil {
		log.Printf("Failed to insert beer into database: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusCreated, newBeer)
}

func addBeerDB(newBeer *beer) error {

	db_host := os.Getenv("DB_HOST")
	db_pass := os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=postgres "+
		"password=%s dbname=postgres sslmode=disable",
		db_host, db_pass)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(newBeer.Name)
	_, err = db.Exec("INSERT INTO beer (name) VALUES ($1);", newBeer.Name)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
