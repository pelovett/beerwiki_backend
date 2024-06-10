package beer

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/db_wrapper"
)

func GetRandomBeer(c *gin.Context) {

	rows, err := db_wrapper.Query("SELECT url_name FROM beer ORDER BY RANDOM() LIMIT 1;")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Beer not found",
		})
		log.Printf("Failed to run query: %s", err)
		return
	}

	defer rows.Close()

	if !rows.Next() {
		log.Println("No rows returned.")
		return
	}

	var urlName string

	err = rows.Scan(&urlName)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"url_name": urlName})

}
