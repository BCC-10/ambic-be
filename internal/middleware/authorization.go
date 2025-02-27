package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Authorization(ctx *fiber.Ctx) error {
	isActive := ctx.Locals("isActive")

	if isActive == false {
		return res.Forbidden(ctx)
	}

	return ctx.Next()
}
