package repository

import (
	"api-server/models"
	"errors"

	"gorm.io/gorm"
)

// MovieRepository defines the contract for the movie repository
type MovieRepository interface {
	FindAll(page, limit int, genreID, directorID *uint, minRating *float64) ([]models.Movie, int64, error)
	FindByID(id uint) (*models.Movie, error)
	Create(movie *models.Movie) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	FindByGenre(genreID uint, page, limit int) ([]models.Movie, int64, error)
	FindByDirector(directorID uint, page, limit int) ([]models.Movie, int64, error)
	FindByActor(actorID uint, page, limit int) ([]models.Movie, int64, error)
	SearchByTitle(title string, page, limit int) ([]models.Movie, int64, error)
	GetTopRated(limit int) ([]models.Movie, error)
}

// gormMovieRepository is the concrete implementation using GORM
type gormMovieRepository struct {
	db *gorm.DB
}

// NewMovieRepository creates a new repository instance with dependency injection
func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &gormMovieRepository{db: db}
}

func (r *gormMovieRepository) FindAll(page, limit int, genreID, directorID *uint, minRating *float64) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.Model(&models.Movie{}).Preload("Genre").Preload("Director").Preload("Actors")

	// Apply filters
	if genreID != nil {
		query = query.Where("genre_id = ?", *genreID)
	}
	if directorID != nil {
		query = query.Where("director_id = ?", *directorID)
	}
	if minRating != nil {
		query = query.Where("rating >= ?", *minRating)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *gormMovieRepository) FindByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	err := r.db.Preload("Genre").Preload("Director").Preload("Actors").Preload("Reviews.User").First(&movie, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}

func (r *gormMovieRepository) Create(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *gormMovieRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.Movie{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}
	return nil
}

func (r *gormMovieRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Movie{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}
	return nil
}

func (r *gormMovieRepository) FindByGenre(genreID uint, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.Model(&models.Movie{}).Where("genre_id = ?", genreID).Preload("Genre").Preload("Director").Preload("Actors")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *gormMovieRepository) FindByDirector(directorID uint, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.Model(&models.Movie{}).Where("director_id = ?", directorID).Preload("Genre").Preload("Director").Preload("Actors")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *gormMovieRepository) FindByActor(actorID uint, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.Model(&models.Movie{}).Joins("JOIN movie_actors ON movies.id = movie_actors.movie_id").
		Where("movie_actors.actor_id = ?", actorID).Preload("Genre").Preload("Director").Preload("Actors")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *gormMovieRepository) SearchByTitle(title string, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.Model(&models.Movie{}).Where("title LIKE ?", "%"+title+"%").
		Preload("Genre").Preload("Director").Preload("Actors")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *gormMovieRepository) GetTopRated(limit int) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Model(&models.Movie{}).Order("rating DESC").Limit(limit).
		Preload("Genre").Preload("Director").Preload("Actors").Find(&movies).Error
	return movies, err
} 