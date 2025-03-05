package fiber

import (
	"ambic/internal/domain/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"time"
)

func New(env *env.Env) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
		BodyLimit:   int(env.MaxUploadSize * 1024 * 1024),
	})

	app.Use(logger.New())

	app.Use(cors.New())

	return app
}
