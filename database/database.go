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
		Logger: logger.Default.LogMode(logger.Info),
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
			{Name: "Acción", Description: "Películas de acción y aventura"},
			{Name: "Comedia", Description: "Películas cómicas y divertidas"},
			{Name: "Drama", Description: "Películas dramáticas"},
			{Name: "Terror", Description: "Películas de terror y suspenso"},
			{Name: "Ciencia Ficción", Description: "Películas de ciencia ficción"},
			{Name: "Romance", Description: "Películas románticas"},
		}
		DB.Create(&genres)

		// Create sample directors
		directors := []models.Director{
			{
				Name:        "Christopher Nolan",
				Biography:   "Director británico conocido por Inception, Interstellar y The Dark Knight",
				Nationality: "Británico",
			},
			{
				Name:        "Quentin Tarantino",
				Biography:   "Director estadounidense conocido por Pulp Fiction y Kill Bill",
				Nationality: "Estadounidense",
			},
			{
				Name:        "Greta Gerwig",
				Biography:   "Directora estadounidense conocida por Lady Bird y Barbie",
				Nationality: "Estadounidense",
			},
		}
		DB.Create(&directors)

		// Create sample actors
		actors := []models.Actor{
			{
				Name:        "Leonardo DiCaprio",
				Biography:   "Actor estadounidense ganador del Oscar",
				Nationality: "Estadounidense",
			},
			{
				Name:        "Margot Robbie",
				Biography:   "Actriz australiana conocida por Barbie y Suicide Squad",
				Nationality: "Australiana",
			},
			{
				Name:        "Tom Hardy",
				Biography:   "Actor británico conocido por Mad Max y Venom",
				Nationality: "Británico",
			},
			{
				Name:        "Emma Stone",
				Biography:   "Actriz estadounidense ganadora del Oscar",
				Nationality: "Estadounidense",
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
				Description: "Un ladrón que roba información corporativa a través del uso de la tecnología de compartir sueños, recibe la tarea inversa de plantar una idea en la mente de un CEO.",
				ReleaseYear: 2010,
				Duration:    148,
				Rating:      8.8,
				PosterURL:   "https://example.com/inception.jpg",
				TrailerURL:  "https://example.com/inception-trailer.mp4",
				GenreID:     &genres[4].ID, // Ciencia Ficción
				DirectorID:  &directors[0].ID,
			},
			{
				Title:       "Barbie",
				Description: "Barbie sufre una crisis existencial y viaja al mundo real para encontrar la verdadera felicidad.",
				ReleaseYear: 2023,
				Duration:    114,
				Rating:      7.0,
				PosterURL:   "https://example.com/barbie.jpg",
				TrailerURL:  "https://example.com/barbie-trailer.mp4",
				GenreID:     &genres[1].ID, // Comedia
				DirectorID:  &directors[2].ID,
			},
			{
				Title:       "Pulp Fiction",
				Description: "Las vidas de dos sicarios de la mafia, un boxeador, la esposa de un gángster y dos bandidos se entrelazan en cuatro historias de violencia y redención.",
				ReleaseYear: 1994,
				Duration:    154,
				Rating:      8.9,
				PosterURL:   "https://example.com/pulp-fiction.jpg",
				TrailerURL:  "https://example.com/pulp-fiction-trailer.mp4",
				GenreID:     &genres[0].ID, // Acción
				DirectorID:  &directors[1].ID,
			},
			{
				Title:       "Poor Things",
				Description: "La increíble evolución de Bella Baxter, una joven traída de vuelta a la vida por el brillante y poco ortodoxo Dr. Godwin Baxter.",
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
				Comment: "Una obra maestra del cine. Nolan nunca decepciona.",
			},
			{
				MovieID: movies[1].ID,
				UserID:  users[1].ID,
				Rating:  7.5,
				Comment: "Divertida y con un mensaje importante. Margot Robbie está increíble.",
			},
			{
				MovieID: movies[2].ID,
				UserID:  users[2].ID,
				Rating:  9.5,
				Comment: "Clásico absoluto. Tarantino en su mejor momento.",
			},
		}
		DB.Create(&reviews)

		log.Println("Database seeded with sample movie data")
	}
} 