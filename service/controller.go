package service

import (
	"api-server/database"
	"api-server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s *MovieService) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	movie, err := s.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movie,
	})
}

func (s *MovieService) Find(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	genreIDStr := c.Query("genre_id")
	directorIDStr := c.Query("director_id")
	minRatingStr := c.Query("min_rating")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse filters
	var genreID, directorID *uint
	var minRating *float64

	if genreIDStr != "" {
		if id, err := strconv.ParseUint(genreIDStr, 10, 32); err == nil {
			genreID = &[]uint{uint(id)}[0]
		}
	}

	if directorIDStr != "" {
		if id, err := strconv.ParseUint(directorIDStr, 10, 32); err == nil {
			directorID = &[]uint{uint(id)}[0]
		}
	}

	if minRatingStr != "" {
		if rating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			minRating = &rating
		}
	}

	movies, total, err := s.repo.FindAll(page, limit, genreID, directorID, minRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *MovieService) Create(c *gin.Context) {
	var req models.MovieCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	movie := models.Movie{
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

	if err := s.repo.Create(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Associate actors if provided
	if len(req.ActorIDs) > 0 {
		var actors []models.Actor
		database.DB.Find(&actors, req.ActorIDs)
		database.DB.Model(&movie).Association("Actors").Append(actors)
	}

	// Get the created movie with relations
	createdMovie, err := s.repo.FindByID(movie.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Movie created but could not retrieve data",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": createdMovie,
	})
}

func (s *MovieService) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	var req models.MovieUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ReleaseYear != nil {
		updates["release_year"] = *req.ReleaseYear
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.Rating != nil {
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

	if err := s.repo.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update actors if provided
	if len(req.ActorIDs) > 0 {
		var movie models.Movie
		database.DB.First(&movie, id)
		var actors []models.Actor
		database.DB.Find(&actors, req.ActorIDs)
		database.DB.Model(&movie).Association("Actors").Replace(actors)
	}

	// Get updated movie
	movie, err := s.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Movie updated but could not retrieve updated data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movie,
	})
}

func (s *MovieService) Remove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	if err := s.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully",
	})
}

func (s *MovieService) Search(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title parameter is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := s.repo.SearchByTitle(title, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *MovieService) TopRated(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	movies, err := s.repo.GetTopRated(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
	})
}

func (s *MovieService) ByGenre(c *gin.Context) {
	genreIDStr := c.Param("id")
	genreID, err := strconv.ParseUint(genreIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid genre ID format",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := s.repo.FindByGenre(uint(genreID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *MovieService) ByDirector(c *gin.Context) {
	directorIDStr := c.Param("id")
	directorID, err := strconv.ParseUint(directorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid director ID format",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := s.repo.FindByDirector(uint(directorID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *MovieService) ByActor(c *gin.Context) {
	actorIDStr := c.Param("id")
	actorID, err := strconv.ParseUint(actorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid actor ID format",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := s.repo.FindByActor(uint(actorID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
