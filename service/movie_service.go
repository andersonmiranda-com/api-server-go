package service

import (
	"api-server/models"
	"api-server/repository"
	"errors"
)

// MovieService defines the contract for movie business logic
type MovieService interface {
	GetMovie(id uint) (*models.Movie, error)
	GetMovies(page, limit int, genreID, directorID *uint, minRating *float64) ([]models.Movie, int64, error)
	CreateMovie(req *models.MovieCreateRequest) (*models.Movie, error)
	UpdateMovie(id uint, req *models.MovieUpdateRequest) (*models.Movie, error)
	DeleteMovie(id uint) error
	SearchMovies(title string, page, limit int) ([]models.Movie, int64, error)
	GetTopRatedMovies(limit int) ([]models.Movie, error)
	GetMoviesByGenre(genreID uint, page, limit int) ([]models.Movie, int64, error)
	GetMoviesByDirector(directorID uint, page, limit int) ([]models.Movie, int64, error)
	GetMoviesByActor(actorID uint, page, limit int) ([]models.Movie, int64, error)
}

// movieServiceImpl is the concrete implementation of the service
type movieServiceImpl struct {
	repo repository.MovieRepository
}

// NewMovieService creates a new service instance with dependency injection
func NewMovieService(repo repository.MovieRepository) MovieService {
	return &movieServiceImpl{repo: repo}
}

func (s *movieServiceImpl) GetMovie(id uint) (*models.Movie, error) {
	if id == 0 {
		return nil, errors.New("invalid movie ID")
	}
	return s.repo.FindByID(id)
}

func (s *movieServiceImpl) GetMovies(page, limit int, genreID, directorID *uint, minRating *float64) ([]models.Movie, int64, error) {
	// Business validations
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	return s.repo.FindAll(page, limit, genreID, directorID, minRating)
}

func (s *movieServiceImpl) CreateMovie(req *models.MovieCreateRequest) (*models.Movie, error) {
	// Business validations
	if req.Title == "" {
		return nil, errors.New("movie title is required")
	}
	if req.ReleaseYear < 1888 || req.ReleaseYear > 2030 {
		return nil, errors.New("invalid release year")
	}
	if req.Duration <= 0 {
		return nil, errors.New("duration must be positive")
	}
	if req.Rating < 0 || req.Rating > 10 {
		return nil, errors.New("rating must be between 0 and 10")
	}

	movie := &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		ReleaseYear: req.ReleaseYear,
		Duration:    req.Duration,
		Rating:      req.Rating,
		PosterURL:   req.PosterURL,
		TrailerURL:  req.TrailerURL,
		GenreID:     req.GenreID,
		DirectorID:  req.DirectorID,
	}

	if err := s.repo.Create(movie); err != nil {
		return nil, err
	}

	// Return the created movie with all its relations
	return s.repo.FindByID(movie.ID)
}

func (s *movieServiceImpl) UpdateMovie(id uint, req *models.MovieUpdateRequest) (*models.Movie, error) {
	if id == 0 {
		return nil, errors.New("invalid movie ID")
	}

	// Verify that the movie exists
	existingMovie, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Build the updates map
	updates := make(map[string]interface{})
	
	if req.Title != nil {
		if *req.Title == "" {
			return nil, errors.New("movie title cannot be empty")
		}
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ReleaseYear != nil {
		if *req.ReleaseYear < 1888 || *req.ReleaseYear > 2030 {
			return nil, errors.New("invalid release year")
		}
		updates["release_year"] = *req.ReleaseYear
	}
	if req.Duration != nil {
		if *req.Duration <= 0 {
			return nil, errors.New("duration must be positive")
		}
		updates["duration"] = *req.Duration
	}
	if req.Rating != nil {
		if *req.Rating < 0 || *req.Rating > 10 {
			return nil, errors.New("rating must be between 0 and 10")
		}
		updates["rating"] = *req.Rating
	}
	if req.PosterURL != nil {
		updates["poster_url"] = *req.PosterURL
	}
	if req.TrailerURL != nil {
		updates["trailer_url"] = *req.TrailerURL
	}
	if req.GenreID != nil {
		updates["genre_id"] = *req.GenreID
	}
	if req.DirectorID != nil {
		updates["director_id"] = *req.DirectorID
	}

	if len(updates) == 0 {
		return existingMovie, nil // No changes
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	// Return the updated movie
	return s.repo.FindByID(id)
}

func (s *movieServiceImpl) DeleteMovie(id uint) error {
	if id == 0 {
		return errors.New("invalid movie ID")
	}
	return s.repo.Delete(id)
}

func (s *movieServiceImpl) SearchMovies(title string, page, limit int) ([]models.Movie, int64, error) {
	if title == "" {
		return nil, 0, errors.New("search title is required")
	}
	
	// Pagination validations
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	return s.repo.SearchByTitle(title, page, limit)
}

func (s *movieServiceImpl) GetTopRatedMovies(limit int) ([]models.Movie, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.GetTopRated(limit)
}

func (s *movieServiceImpl) GetMoviesByGenre(genreID uint, page, limit int) ([]models.Movie, int64, error) {
	if genreID == 0 {
		return nil, 0, errors.New("invalid genre ID")
	}
	
	// Pagination validations
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	return s.repo.FindByGenre(genreID, page, limit)
}

func (s *movieServiceImpl) GetMoviesByDirector(directorID uint, page, limit int) ([]models.Movie, int64, error) {
	if directorID == 0 {
		return nil, 0, errors.New("invalid director ID")
	}
	
	// Pagination validations
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	return s.repo.FindByDirector(directorID, page, limit)
}

func (s *movieServiceImpl) GetMoviesByActor(actorID uint, page, limit int) ([]models.Movie, int64, error) {
	if actorID == 0 {
		return nil, 0, errors.New("invalid actor ID")
	}
	
	// Pagination validations
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	return s.repo.FindByActor(actorID, page, limit)
} 