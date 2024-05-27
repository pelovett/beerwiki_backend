package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/api/handlers"
	"github.com/pelovett/beerwiki_backend/middleware"
)

func main() {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())

	// Debug
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Account Management
	r.POST("user/create-account", handlers.CreateUser)
	r.POST("user/login", handlers.LoginUser)
	r.POST("user/verify", handlers.VerifyUser)

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
