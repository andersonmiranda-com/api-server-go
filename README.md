# Movie API Server - General API Concept

A modern REST API built with Go, Gin and SQLite to manage a movie database with genres, directors, actors and reviews.

## 🎬 Features

- **Framework**: Gin (Go)
- **Database**: SQLite with GORM
- **Validation**: go-playground/validator
- **Entities**: Movies, Genres, Directors, Actors, Users, Reviews
- **Relationships**: Many-to-many between movies and actors
- **Sample Data**: Includes popular movies with complete data

## 📋 Available Endpoints

### Health Check
- `GET /health` - API status

### Movies
- `GET /movies` - List movies (with filters and pagination)
- `GET /movies/:id` - Get movie by ID
- `POST /movies` - Create new movie
- `PUT /movies/:id` - Update movie
- `DELETE /movies/:id` - Delete movie
- `GET /movies/search?title=inception` - Search movies by title
- `GET /movies/top-rated?limit=10` - Top rated movies

### Genres
- `GET /genres/:id/movies` - Movies by genre

### Directors
- `GET /directors/:id/movies` - Movies by director

### Actors
- `GET /actors/:id/movies` - Movies by actor

### Query Parameters
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `genre_id` - Filter by genre
- `director_id` - Filter by director
- `min_rating` - Filter by minimum rating
- `title` - Search by title (for search endpoint)

## 🛠️ Installation and Usage

### Prerequisites
- Go 1.21 or higher

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd api-server

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

The API will be available at `http://localhost:4444`

### Sample Data
The application includes sample data:
- **6 genres**: Action, Comedy, Drama, Horror, Science Fiction, Romance
- **3 directors**: Christopher Nolan, Quentin Tarantino, Greta Gerwig
- **4 actors**: Leonardo DiCaprio, Margot Robbie, Tom Hardy, Emma Stone
- **4 movies**: Inception, Barbie, Pulp Fiction, Poor Things
- **3 users** with reviews

## 📝 Usage Examples

### List all movies
```bash
curl http://localhost:4444/movies
```

### Search movies by title
```bash
curl "http://localhost:4444/movies/search?title=inception"
```

### Filter movies by genre
```bash
curl "http://localhost:4444/movies?genre_id=1&min_rating=8.0"
```

### Get top rated movies
```bash
curl "http://localhost:4444/movies/top-rated?limit=5"
```

### Create a new movie
```bash
curl -X POST http://localhost:4444/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "description": "A programmer discovers that reality is a simulation",
    "release_year": 1999,
    "duration": 136,
    "rating": 8.7,
    "poster_url": "https://example.com/matrix.jpg",
    "trailer_url": "https://example.com/matrix-trailer.mp4",
    "genre_id": 5,
    "director_id": 1,
    "actor_ids": [1, 3]
  }'
```

### Get movie by ID
```bash
curl http://localhost:4444/movies/1
```

### Movies by director
```bash
curl http://localhost:4444/directors/1/movies
```

## 🗄️ Database Structure

### Movies
- `id` - Unique ID
- `title` - Movie title
- `description` - Description
- `release_year` - Release year
- `duration` - Duration in minutes
- `rating` - Rating (0-10)
- `poster_url` - Poster URL
- `trailer_url` - Trailer URL
- `genre_id` - Genre ID
- `director_id` - Director ID
- `created_at` - Creation date
- `updated_at` - Update date

### Genres
- `id` - Unique ID
- `name` - Genre name
- `description` - Description
- `created_at` - Creation date
- `updated_at` - Update date

### Directors
- `id` - Unique ID
- `name` - Director name
- `biography` - Biography
- `birth_date` - Birth date
- `nationality` - Nationality
- `created_at` - Creation date
- `updated_at` - Update date

### Actors
- `id` - Unique ID
- `name` - Actor name
- `biography` - Biography
- `birth_date` - Birth date
- `nationality` - Nationality
- `created_at` - Creation date
- `updated_at` - Update date

### Reviews
- `id` - Unique ID
- `movie_id` - Movie ID
- `user_id` - User ID
- `rating` - Rating (1-10)
- `comment` - Comment
- `created_at` - Creation date
- `updated_at` - Update date

### Users
- `id` - Unique ID
- `username` - Username
- `email` - Email
- `created_at` - Creation date
- `updated_at` - Update date

## 🔧 Configuration

The application uses SQLite by default. The database file is automatically created as `api_server.db` in the root directory.

## 🧪 Future Improvements

- [ ] JWT Authentication
- [ ] Watchlist system
- [ ] Recommendations based on preferences
- [ ] Swagger documentation
- [ ] Unit tests
- [ ] Redis cache
- [ ] Advanced search (by actor, director)
- [ ] File upload (posters, trailers)
- [ ] Structured logging
- [ ] Metrics and monitoring
- [ ] Complete CRUD API for genres, directors and actors

## 📄 License

This project is for educational and demonstration purposes.


