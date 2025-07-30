package service

import (
	"api-server/repository"
	"github.com/gin-gonic/gin"
)

type MovieService struct {
	repo *repository.MovieRepository
}

type Service interface {
	Get(c *gin.Context)
	Find(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Remove(c *gin.Context)
	Search(c *gin.Context)
	TopRated(c *gin.Context)
	ByGenre(c *gin.Context)
	ByDirector(c *gin.Context)
	ByActor(c *gin.Context)
}

func NewService() Service {
	return &MovieService{
		repo: repository.NewMovieRepository(),
	}
}
