package database

import (
	"api-server/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("api_server.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Solo mostrar errores
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Genre{}, &models.Director{}, &models.Actor{}, &models.User{}, &models.Movie{}, &models.Review{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data if database is empty
	seedData()
}

func seedData() {
	var count int64
	DB.Model(&models.Genre{}).Count(&count)
	if count == 0 {
		// Create sample genres
		genres := []models.Genre{
			{Name: "Action", Description: "Action and adventure movies"},
			{Name: "Comedy", Description: "Funny and entertaining movies"},
			{Name: "Drama", Description: "Dramatic movies"},
			{Name: "Horror", Description: "Horror and suspense movies"},
			{Name: "Science Fiction", Description: "Science fiction movies"},
			{Name: "Romance", Description: "Romantic movies"},
		}
		DB.Create(&genres)

		// Create sample directors
		directors := []models.Director{
			{
				Name:        "Christopher Nolan",
				Biography:   "British director known for Inception, Interstellar and The Dark Knight",
				Nationality: "British",
			},
			{
				Name:        "Quentin Tarantino",
				Biography:   "American director known for Pulp Fiction and Kill Bill",
				Nationality: "American",
			},
			{
				Name:        "Greta Gerwig",
				Biography:   "American director known for Lady Bird and Barbie",
				Nationality: "American",
			},
		}
		DB.Create(&directors)

		// Create sample actors
		actors := []models.Actor{
			{
				Name:        "Leonardo DiCaprio",
				Biography:   "American actor and Oscar winner",
				Nationality: "American",
			},
			{
				Name:        "Margot Robbie",
				Biography:   "Australian actress known for Barbie and Suicide Squad",
				Nationality: "Australian",
			},
			{
				Name:        "Tom Hardy",
				Biography:   "British actor known for Mad Max and Venom",
				Nationality: "British",
			},
			{
				Name:        "Emma Stone",
				Biography:   "American actress and Oscar winner",
				Nationality: "American",
			},
		}
		DB.Create(&actors)

		// Create sample users
		users := []models.User{
			{Username: "movie_lover", Email: "lover@movies.com"},
			{Username: "cinema_fan", Email: "fan@cinema.com"},
			{Username: "film_critic", Email: "critic@films.com"},
		}
		DB.Create(&users)

		// Create sample movies
		movies := []models.Movie{
			{
				Title:       "Inception",
				Description: "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.",
				ReleaseYear: 2010,
				Duration:    148,
				Rating:      8.8,
				PosterURL:   "https://example.com/inception.jpg",
				TrailerURL:  "https://example.com/inception-trailer.mp4",
				GenreID:     &genres[4].ID, // Science Fiction
				DirectorID:  &directors[0].ID,
			},
			{
				Title:       "Barbie",
				Description: "Barbie suffers an existential crisis and travels to the real world to find true happiness.",
				ReleaseYear: 2023,
				Duration:    114,
				Rating:      7.0,
				PosterURL:   "https://example.com/barbie.jpg",
				TrailerURL:  "https://example.com/barbie-trailer.mp4",
				GenreID:     &genres[1].ID, // Comedy
				DirectorID:  &directors[2].ID,
			},
			{
				Title:       "Pulp Fiction",
				Description: "The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.",
				ReleaseYear: 1994,
				Duration:    154,
				Rating:      8.9,
				PosterURL:   "https://example.com/pulp-fiction.jpg",
				TrailerURL:  "https://example.com/pulp-fiction-trailer.mp4",
				GenreID:     &genres[0].ID, // Action
				DirectorID:  &directors[1].ID,
			},
			{
				Title:       "Poor Things",
				Description: "The incredible evolution of Bella Baxter, a young woman brought back to life by the brilliant and unorthodox scientist Dr. Godwin Baxter.",
				ReleaseYear: 2023,
				Duration:    141,
				Rating:      8.4,
				PosterURL:   "https://example.com/poor-things.jpg",
				TrailerURL:  "https://example.com/poor-things-trailer.mp4",
				GenreID:     &genres[2].ID, // Drama
				DirectorID:  &directors[2].ID,
			},
		}
		DB.Create(&movies)

		// Associate actors with movies
		DB.Model(&movies[0]).Association("Actors").Append(&actors[0], &actors[2]) // Inception
		DB.Model(&movies[1]).Association("Actors").Append(&actors[1])             // Barbie
		DB.Model(&movies[2]).Association("Actors").Append(&actors[0])             // Pulp Fiction
		DB.Model(&movies[3]).Association("Actors").Append(&actors[3])             // Poor Things

		// Create sample reviews
		reviews := []models.Review{
			{
				MovieID: movies[0].ID,
				UserID:  users[0].ID,
				Rating:  9.0,
				Comment: "A masterpiece of cinema. Nolan never disappoints.",
			},
			{
				MovieID: movies[1].ID,
				UserID:  users[1].ID,
				Rating:  7.5,
				Comment: "Funny and with an important message. Margot Robbie is incredible.",
			},
			{
				MovieID: movies[2].ID,
				UserID:  users[2].ID,
				Rating:  9.5,
				Comment: "Absolute classic. Tarantino at his best.",
			},
		}
		DB.Create(&reviews)

		log.Println("Database seeded with sample movie data")
	}
} 