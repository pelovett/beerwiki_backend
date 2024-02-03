package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/api/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db_host := os.Getenv("DB_HOST")
	db_pass := os.Getenv("DB_PASSWORD")
	psqlInfo := fmt.Sprintf("host=%s port=5432 user=postgres "+
		"password=%s dbname=postgres sslmode=disable",
		db_host, db_pass)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	// fmt.Printf("the host is %s\n", db_host)
	// fmt.Printf("the pass is %s\n", db_pass)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		result, err := db.Query("SELECT name FROM beer;")
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Very bad",
			})
			return
		}
		var beerName string
		if result.Next() {
			fmt.Printf("Result next ran\n")
			if err := result.Scan(&beerName); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("The result is %s\n", beerName)

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/beer/:id", handlers.GetBeer)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
