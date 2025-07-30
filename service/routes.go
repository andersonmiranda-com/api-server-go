package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine, cts Service) {

	// Health check
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "ok",
			"status": http.StatusOK,
		})
	})

	// Movies routes
	movies := app.Group("/movies")
	movies.GET("/", cts.Find)                    // GET /movies?page=1&limit=10&genre_id=1&director_id=1&min_rating=8.0
	movies.GET("/search", cts.Search)            // GET /movies/search?title=inception
	movies.GET("/top-rated", cts.TopRated)       // GET /movies/top-rated?limit=10
	movies.GET("/:id", cts.Get)                  // GET /movies/1
	movies.POST("/", cts.Create)                 // POST /movies
	movies.PUT("/:id", cts.Update)               // PUT /movies/1
	movies.DELETE("/:id", cts.Remove)            // DELETE /movies/1

	// Genre routes
	genres := app.Group("/genres")
	genres.GET("/:id/movies", cts.ByGenre)       // GET /genres/1/movies

	// Director routes
	directors := app.Group("/directors")
	directors.GET("/:id/movies", cts.ByDirector) // GET /directors/1/movies

	// Actor routes
	actors := app.Group("/actors")
	actors.GET("/:id/movies", cts.ByActor)       // GET /actors/1/movies

}
