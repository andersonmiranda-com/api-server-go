package main

import (
	"api-server/database"
	"api-server/handler"
	"api-server/repository"
	"api-server/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create Gin app
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()

	// Middleware
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Dependency Injection - Hexagonal Architecture
	// 1. Create repository (data layer)
	movieRepo := repository.NewMovieRepository(database.DB)
	
	// 2. Create service (business logic)
	movieService := service.NewMovieService(movieRepo)
	
	// 3. Create handler (HTTP adapter)
	movieHandler := handler.NewMovieHandler(movieService)

	// 4. Configure routes
	handler.SetupRoutes(app, movieHandler)

	// Start server
	log.Println("🚀 Server starting on port 4444...")
	log.Fatal(app.Run(":4444"))
}
