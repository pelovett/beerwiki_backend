package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/api/handlers"
)

func main() {
	r := gin.Default()

	// Debug
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Account Management
	r.POST("/create-account", handlers.CreateUser)
	r.GET("/login", handlers.VerifyUser)

	// Beer
	r.GET("/beer/:id", handlers.GetBeer)
	r.GET("/beer/name/:name", handlers.GetBeerByUrlName)
	r.POST("/beer", handlers.PostBeer)

	// Random
	r.GET("/randombeer/", handlers.GetRandomBeer)

	// Search
	r.GET("/search/beer", handlers.SearchBeerByName)
	r.Run()
}
