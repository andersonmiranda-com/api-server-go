package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" validate:"required"`
	Description string         `json:"description"`
	ReleaseYear int            `json:"release_year" validate:"min=1888,max=2030"`
	Duration    int            `json:"duration" validate:"min=1"` // in minutes
	Rating      float64        `json:"rating" validate:"min=0,max=10"`
	PosterURL   string         `json:"poster_url"`
	TrailerURL  string         `json:"trailer_url"`
	GenreID     *uint          `json:"genre_id"`
	Genre       *Genre         `json:"genre,omitempty"`
	DirectorID  *uint          `json:"director_id"`
	Director    *Director      `json:"director,omitempty"`
	Actors      []Actor        `json:"actors,omitempty" gorm:"many2many:movie_actors;"`
	Reviews     []Review       `json:"reviews,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Genre struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Movies      []Movie        `json:"movies,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Director struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" validate:"required"`
	Biography   string         `json:"biography"`
	BirthDate   *time.Time     `json:"birth_date"`
	Nationality string         `json:"nationality"`
	Movies      []Movie        `json:"movies,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Actor struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" validate:"required"`
	Biography   string         `json:"biography"`
	BirthDate   *time.Time     `json:"birth_date"`
	Nationality string         `json:"nationality"`
	Movies      []Movie        `json:"movies,omitempty" gorm:"many2many:movie_actors;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MovieID   uint           `json:"movie_id"`
	Movie     Movie          `json:"movie,omitempty"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user,omitempty"`
	Rating    float64        `json:"rating" validate:"min=1,max=10"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" validate:"required"`
	Email     string         `json:"email" validate:"required,email"`
	Reviews   []Review       `json:"reviews,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Request/Response models
type MovieCreateRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	ReleaseYear int     `json:"release_year" validate:"min=1888,max=2030"`
	Duration    int     `json:"duration" validate:"min=1"`
	Rating      float64 `json:"rating" validate:"min=0,max=10"`
	PosterURL   string  `json:"poster_url"`
	TrailerURL  string  `json:"trailer_url"`
	GenreID     *uint   `json:"genre_id"`
	DirectorID  *uint   `json:"director_id"`
	ActorIDs    []uint  `json:"actor_ids"`
}

type MovieUpdateRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	ReleaseYear *int     `json:"release_year" validate:"omitempty,min=1888,max=2030"`
	Duration    *int     `json:"duration" validate:"omitempty,min=1"`
	Rating      *float64 `json:"rating" validate:"omitempty,min=0,max=10"`
	PosterURL   *string  `json:"poster_url"`
	TrailerURL  *string  `json:"trailer_url"`
	GenreID     *uint    `json:"genre_id"`
	DirectorID  *uint    `json:"director_id"`
	ActorIDs    []uint   `json:"actor_ids"`
}

type ReviewCreateRequest struct {
	MovieID uint    `json:"movie_id" validate:"required"`
	UserID  uint    `json:"user_id" validate:"required"`
	Rating  float64 `json:"rating" validate:"min=1,max=10"`
	Comment string  `json:"comment"`
} 