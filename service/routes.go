package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cts Service) {

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).
			JSON(map[string]interface{}{
				"health": "ok",
				"status": http.StatusOK,
			})
	})

	// Movies routes
	movies := app.Group("/movies")
	movies.Get("/", cts.Find)                    // GET /movies?page=1&limit=10&genre_id=1&director_id=1&min_rating=8.0
	movies.Get("/search", cts.Search)            // GET /movies/search?title=inception
	movies.Get("/top-rated", cts.TopRated)       // GET /movies/top-rated?limit=10
	movies.Get("/:id", cts.Get)                  // GET /movies/1
	movies.Post("/", cts.Create)                 // POST /movies
	movies.Put("/:id", cts.Update)               // PUT /movies/1
	movies.Delete("/:id", cts.Remove)            // DELETE /movies/1

	// Genre routes
	genres := app.Group("/genres")
	genres.Get("/:id/movies", cts.ByGenre)       // GET /genres/1/movies

	// Director routes
	directors := app.Group("/directors")
	directors.Get("/:id/movies", cts.ByDirector) // GET /directors/1/movies

	// Actor routes
	actors := app.Group("/actors")
	actors.Get("/:id/movies", cts.ByActor)       // GET /actors/1/movies

}
