package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
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


	log.Println(newBeer.Name)
    _, err := db_wrapper.ExecSQL("INSERT INTO beer (name) VALUES ($1);", newBeer.Name)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
