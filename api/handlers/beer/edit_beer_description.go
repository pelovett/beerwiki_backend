package beer

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

type descriptionChangeRequest struct {
	URLName   string    `json:"url_name"`
	PageIPAML string `json:"page_ipa_ml"`
}

func PostBeerDescription(c *gin.Context) {
	var newDescription descriptionChangeRequest;
	log.Printf("Posting beer description");

	if err := c.BindJSON(&newDescription); err != nil {
		log.Printf("Failed to bind inputted json to description change request with err %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := changeDescriptionDB(&newDescription); err != nil {
		log.Printf("Failed to change beer description: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusCreated, newDescription)
}

func changeDescriptionDB(change *descriptionChangeRequest) error {
	_, err := db_wrapper.Exec("UPDATE beer SET page_ipa_ml=$1 WHERE url_name = $2;",
	    change.PageIPAML, change.URLName)
	
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
