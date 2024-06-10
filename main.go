package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/api/handlers/beer"
	"github.com/pelovett/beerwiki_backend/api/handlers/image"
	"github.com/pelovett/beerwiki_backend/api/handlers/user"
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
	r.POST("user/create-account", user.CreateUser)
	r.POST("user/login", user.LoginUser)
	r.POST("user/verify", user.VerifyUser)

	// Beer
	r.GET("/beer/:id", beer.GetBeer)
	r.GET("/beer/name/:name", beer.GetBeerByUrlName)
	r.POST("/beer", beer.PostBeer)

	// Image
	r.GET("/image/upload", image.GetImageUploadURL)
	r.POST("/image/upload/complete", image.PostImageUploadComplete)

	// Random
	r.GET("/randombeer/", beer.GetRandomBeer)

	// Search
	r.GET("/search/beer", beer.SearchBeerByName)

	r.Run()
}
