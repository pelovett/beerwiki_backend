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
	id   int    `json:"id"`
	name string `json:"name"`
}

func PostBeer(c *gin.Context) {

	var newBeer beer

	if err := c.BindJSON(&newBeer); err != nil {
		log.Println(`Failed to bind inputted json to beer with err${err}`)
		return
	}

	log.Println(newBeer)

	c.IndentedJSON(http.StatusCreated, newBeer)
	// c.JSON(http.StatusCreated, gin.H{
	//  "message": newBeer.name,
	// })

	fmt.Println(newBeer)

	if err := addBeerDB(&newBeer); err != nil {
		log.Println("Failed to insert beer into database")
		log.Println(err)
	}
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

	log.Println(newBeer.name)
	_, err = db.Exec("INSERT INTO beer (name) VALUES ($1);", newBeer.name)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
