package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) EnsurePartner(ctx *fiber.Ctx) error {
	isPartner := ctx.Locals("isPartner").(bool)
	if !isPartner {
		return res.Forbidden(ctx)
	}

	return ctx.Next()
}
