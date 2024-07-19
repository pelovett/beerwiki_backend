package beer

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

type descriptionChangeRequest struct {
	ID        int    `json:"id"`
	PageIPAML string `json:"page_ipa_ml"`
}

func PostBeerDescription(c *gin.Context) {
	var newDescription descriptionChangeRequest;

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

}

func changeDescriptionDB(change *descriptionChangeRequest) error {
	_, err := db_wrapper.Exec("UPDATE beer SET page_ipa_ml=$1 WHERE id = $2;",
	    change.PageIPAML, change.ID)
	
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
