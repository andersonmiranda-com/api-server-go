package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return Response("testdata", 200, nil, c)
	})

	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(nil))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	app.Get("/error", func(c *fiber.Ctx) error {
		return Response(nil, 500, fmt.Errorf("test error"), c)
	})

	req, _ = http.NewRequest("GET", "/error", bytes.NewBuffer(nil))
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
