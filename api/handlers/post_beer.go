package handlers

import (
	"log"
	"net/http"
    "fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

type beer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	URLName   string `json:"url_name"`
	PageIPAML string `json:"page_ipa_ml"`
}

func PostBeer(c *gin.Context) {

	var newBeer beer

	if err := c.BindJSON(&newBeer); err != nil {
		log.Printf("Failed to bind inputted json to beer with err %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	legal_url_regex, err := regexp.Compile(`^[A-Za-z0-9\(\)\_]*$`)
	if err != nil {
		log.Printf("Failed to compile Regex??? Shouldn't happen")
		c.Status(http.StatusInternalServerError)
		return
	}

	if !legal_url_regex.Match([]byte(newBeer.URLName)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Illegal url_name: %s", newBeer.URLName),
		})
		return
	}

	if err := addBeerDB(&newBeer); err != nil {
		log.Printf("Failed to insert beer into database: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusCreated, newBeer)
}

func addBeerDB(newBeer *beer) error {
	_, err := db_wrapper.ExecSQL("INSERT INTO beer (name, url_name, page_ipa_ml) VALUES ($1, $2 ,$3);",
		newBeer.Name, newBeer.URLName, newBeer.PageIPAML)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
