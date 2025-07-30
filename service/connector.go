package service

import (
	"api-server/repository"
	"github.com/gofiber/fiber/v2"
)

type MovieService struct {
	repo *repository.MovieRepository
}

type Service interface {
	Get(c *fiber.Ctx) error
	Find(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Remove(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	TopRated(c *fiber.Ctx) error
	ByGenre(c *fiber.Ctx) error
	ByDirector(c *fiber.Ctx) error
	ByActor(c *fiber.Ctx) error
}

func NewService() Service {
	return &MovieService{
		repo: repository.NewMovieRepository(),
	}
}
