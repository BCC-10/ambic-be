package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) EnsureNotPartner(ctx *fiber.Ctx) error {
	isPartner := ctx.Locals("isPartner").(bool)
	if isPartner {
		return res.Unauthorized(ctx)
	}

	return ctx.Next()
}
