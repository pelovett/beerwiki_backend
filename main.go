package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pelovett/beerwiki_backend/api/handlers/beer"
	"github.com/pelovett/beerwiki_backend/api/handlers/image"
	"github.com/pelovett/beerwiki_backend/api/handlers/user"
	"github.com/pelovett/beerwiki_backend/api/handlers/util"
	"github.com/pelovett/beerwiki_backend/middleware"
)

func main() {
	r := gin.Default()
	r.Use(gin.LoggerWithWriter(os.Stdout))

	//  Load assets into memory
	cityList, err := util.GetCityList()
	if err != nil {
		log.Fatal(err)
	}

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
		userRoutes.GET("/logout-user", user.LogoutUser)

		userRoutes.POST("/create-account", user.CreateUser)
		userRoutes.POST("/confirmation", user.ConfirmUser)
		userRoutes.POST("/forgot-password", user.ForgotPassword)
		userRoutes.POST("/change-password", user.ChangePassword)
		userRoutes.POST("/login", user.LoginUser)

		userRoutes.Use(middleware.Login())
		userRoutes.POST("/verify", user.VerifyUser)
	}

	// Beer
	beerRoutes := r.Group("/beer")
	{
		beerRoutes.GET("/:id", beer.GetBeer)
		beerRoutes.GET("/name/:name", beer.GetBeerByUrlName)
		beerRoutes.GET("/random", beer.GetRandomBeer)

		beerRoutes.Use(middleware.Login())
		beerRoutes.POST("/description", beer.PostBeerDescription)
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

	// Utils
	utilRoutes := r.Group("/util")
	{
		utilRoutes.GET("/city-search", util.SearchCities(cityList))
	}

	// Search
	r.GET("/search/beer", beer.SearchBeerByName)

	r.Run()
}
