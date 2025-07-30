package service

import (
	"api-server/models"
	"errors"
	"testing"
)

// MockMovieRepository is a mock implementation of the repository for testing
type MockMovieRepository struct {
	movies map[uint]*models.Movie
}

func NewMockMovieRepository() *MockMovieRepository {
	return &MockMovieRepository{
		movies: make(map[uint]*models.Movie),
	}
}

func (m *MockMovieRepository) FindAll(page, limit int, genreID, directorID *uint, minRating *float64) ([]models.Movie, int64, error) {
	// Simple implementation for testing
	var movies []models.Movie
	for _, movie := range m.movies {
		movies = append(movies, *movie)
	}
	return movies, int64(len(movies)), nil
}

func (m *MockMovieRepository) FindByID(id uint) (*models.Movie, error) {
	if movie, exists := m.movies[id]; exists {
		return movie, nil
	}
	return nil, errors.New("movie not found")
}

func (m *MockMovieRepository) Create(movie *models.Movie) error {
	movie.ID = uint(len(m.movies) + 1)
	m.movies[movie.ID] = movie
	return nil
}

func (m *MockMovieRepository) Update(id uint, updates map[string]interface{}) error {
	if _, exists := m.movies[id]; !exists {
		return errors.New("movie not found")
	}
	// Simple implementation for testing
	return nil
}

func (m *MockMovieRepository) Delete(id uint) error {
	if _, exists := m.movies[id]; !exists {
		return errors.New("movie not found")
	}
	delete(m.movies, id)
	return nil
}

func (m *MockMovieRepository) FindByGenre(genreID uint, page, limit int) ([]models.Movie, int64, error) {
	return []models.Movie{}, 0, nil
}

func (m *MockMovieRepository) FindByDirector(directorID uint, page, limit int) ([]models.Movie, int64, error) {
	return []models.Movie{}, 0, nil
}

func (m *MockMovieRepository) FindByActor(actorID uint, page, limit int) ([]models.Movie, int64, error) {
	return []models.Movie{}, 0, nil
}

func (m *MockMovieRepository) SearchByTitle(title string, page, limit int) ([]models.Movie, int64, error) {
	return []models.Movie{}, 0, nil
}

func (m *MockMovieRepository) GetTopRated(limit int) ([]models.Movie, error) {
	return []models.Movie{}, nil
}

// TestGetMovie tests the GetMovie method of the service
func TestGetMovie(t *testing.T) {
	// Arrange
	mockRepo := NewMockMovieRepository()
	service := NewMovieService(mockRepo)
	
	// Add a test movie
	testMovie := &models.Movie{
		ID:    1,
		Title: "Test Movie",
	}
	mockRepo.movies[1] = testMovie

	// Act
	movie, err := service.GetMovie(1)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if movie == nil {
		t.Error("Expected movie, got nil")
	}
	if movie.Title != "Test Movie" {
		t.Errorf("Expected 'Test Movie', got %s", movie.Title)
	}
}

// TestGetMovieNotFound tests the case when the movie doesn't exist
func TestGetMovieNotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockMovieRepository()
	service := NewMovieService(mockRepo)

	// Act
	movie, err := service.GetMovie(999)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if movie != nil {
		t.Error("Expected nil movie, got movie")
	}
	if err.Error() != "movie not found" {
		t.Errorf("Expected 'movie not found', got %s", err.Error())
	}
}

// TestGetMovieInvalidID tests the case with invalid ID
func TestGetMovieInvalidID(t *testing.T) {
	// Arrange
	mockRepo := NewMockMovieRepository()
	service := NewMovieService(mockRepo)

	// Act
	movie, err := service.GetMovie(0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if movie != nil {
		t.Error("Expected nil movie, got movie")
	}
	if err.Error() != "invalid movie ID" {
		t.Errorf("Expected 'invalid movie ID', got %s", err.Error())
	}
}

// TestCreateMovie tests movie creation
func TestCreateMovie(t *testing.T) {
	// Arrange
	mockRepo := NewMockMovieRepository()
	service := NewMovieService(mockRepo)
	
	req := &models.MovieCreateRequest{
		Title:       "New Movie",
		Description: "A test movie",
		ReleaseYear: 2023,
		Duration:    120,
		Rating:      8.5,
	}

	// Act
	movie, err := service.CreateMovie(req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if movie == nil {
		t.Error("Expected movie, got nil")
	}
	if movie.Title != "New Movie" {
		t.Errorf("Expected 'New Movie', got %s", movie.Title)
	}
	if movie.ReleaseYear != 2023 {
		t.Errorf("Expected 2023, got %d", movie.ReleaseYear)
	}
}

// TestCreateMovieInvalidData tests creation with invalid data
func TestCreateMovieInvalidData(t *testing.T) {
	// Arrange
	mockRepo := NewMockMovieRepository()
	service := NewMovieService(mockRepo)
	
	req := &models.MovieCreateRequest{
		Title:       "", // Empty title
		Description: "A test movie",
		ReleaseYear: 2023,
		Duration:    120,
		Rating:      8.5,
	}

	// Act
	movie, err := service.CreateMovie(req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if movie != nil {
		t.Error("Expected nil movie, got movie")
	}
	if err.Error() != "movie title is required" {
		t.Errorf("Expected 'movie title is required', got %s", err.Error())
	}
} 