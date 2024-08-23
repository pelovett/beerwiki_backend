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
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/create-account", user.CreateUser)

		userRoutes.Use(middleware.Login())
		userRoutes.POST("/login", user.LoginUser)
		userRoutes.GET("/verify", user.VerifyUser)
	}

	// Beer
	beerRoutes := r.Group("/beer")
	{
		beerRoutes.GET("/:id", beer.GetBeer)
		beerRoutes.GET("/name/:name", beer.GetBeerByUrlName)
		beerRoutes.GET("/random", beer.GetRandomBeer)

		beerRoutes.Use(middleware.Login())
		beerRoutes.POST("/", beer.PostBeer)
	}

	// Image
	imageRoutes := r.Group("/image")
	{
		imageRoutes.GET("/view/:name", image.GetImageURL)

		//imageRoutes.Use(middleware.Login())
		imageRoutes.GET("/upload", image.GetImageUploadURL)
		imageRoutes.POST("/upload/complete", image.PostImageUploadComplete)
	}

	// Search
	r.GET("/search/beer", beer.SearchBeerByName)

	r.Run()
}
