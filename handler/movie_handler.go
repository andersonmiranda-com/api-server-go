package handler

import (
	"api-server/models"
	"api-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MovieHandler handles HTTP requests related to movies
type MovieHandler struct {
	service service.MovieService
}

// NewMovieHandler creates a new handler instance with dependency injection
func NewMovieHandler(s service.MovieService) *MovieHandler {
	return &MovieHandler{service: s}
}

// Get handles GET /movies/:id
func (h *MovieHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	movie, err := h.service.GetMovie(uint(id))
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

// Find handles GET /movies
func (h *MovieHandler) Find(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	genreIDStr := c.Query("genre_id")
	directorIDStr := c.Query("director_id")
	minRatingStr := c.Query("min_rating")

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

	movies, total, err := h.service.GetMovies(page, limit, genreID, directorID, minRating)
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

// Create handles POST /movies
func (h *MovieHandler) Create(c *gin.Context) {
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

	movie, err := h.service.CreateMovie(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": movie,
	})
}

// Update handles PUT /movies/:id
func (h *MovieHandler) Update(c *gin.Context) {
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

	movie, err := h.service.UpdateMovie(uint(id), &req)
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

// Remove handles DELETE /movies/:id
func (h *MovieHandler) Remove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	if err := h.service.DeleteMovie(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully",
	})
}

// Search handles GET /movies/search
func (h *MovieHandler) Search(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title parameter is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, total, err := h.service.SearchMovies(title, page, limit)
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

// TopRated handles GET /movies/top-rated
func (h *MovieHandler) TopRated(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, err := h.service.GetTopRatedMovies(limit)
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

// ByGenre handles GET /genres/:id/movies
func (h *MovieHandler) ByGenre(c *gin.Context) {
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

	movies, total, err := h.service.GetMoviesByGenre(uint(genreID), page, limit)
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

// ByDirector handles GET /directors/:id/movies
func (h *MovieHandler) ByDirector(c *gin.Context) {
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

	movies, total, err := h.service.GetMoviesByDirector(uint(directorID), page, limit)
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

// ByActor handles GET /actors/:id/movies
func (h *MovieHandler) ByActor(c *gin.Context) {
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

	movies, total, err := h.service.GetMoviesByActor(uint(actorID), page, limit)
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