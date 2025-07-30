package utils

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Response(data interface{}, httpStatus int, err error, c *fiber.Ctx) error {
	if err != nil {
		// Handle different types of errors
		if errors.Is(err, errors.New("not found")) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	} else {
		if data != nil {
			return c.Status(httpStatus).JSON(data)
		} else {
			c.Status(httpStatus)
			return nil
		}
	}
}
