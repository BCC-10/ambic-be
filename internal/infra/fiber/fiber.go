package fiber

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	return app
}
