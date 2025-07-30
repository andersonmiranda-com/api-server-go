package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *gin.Engine, movieHandler *MovieHandler) {
	// Health check
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "ok",
			"status": http.StatusOK,
		})
	})

	// Movies routes
	movies := app.Group("/movies")
	movies.GET("/", movieHandler.Find)                    // GET /movies?page=1&limit=10&genre_id=1&director_id=1&min_rating=8.0
	movies.GET("/search", movieHandler.Search)            // GET /movies/search?title=inception
	movies.GET("/top-rated", movieHandler.TopRated)       // GET /movies/top-rated?limit=10
	movies.GET("/:id", movieHandler.Get)                  // GET /movies/1
	movies.POST("/", movieHandler.Create)                 // POST /movies
	movies.PUT("/:id", movieHandler.Update)               // PUT /movies/1
	movies.DELETE("/:id", movieHandler.Remove)            // DELETE /movies/1

	// Genre routes
	genres := app.Group("/genres")
	genres.GET("/:id/movies", movieHandler.ByGenre)       // GET /genres/1/movies

	// Director routes
	directors := app.Group("/directors")
	directors.GET("/:id/movies", movieHandler.ByDirector) // GET /directors/1/movies

	// Actor routes
	actors := app.Group("/actors")
	actors.GET("/:id/movies", movieHandler.ByActor)       // GET /actors/1/movies
} 