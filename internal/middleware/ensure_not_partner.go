package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) EnsureNotPartner(ctx *fiber.Ctx) error {
	partnerId := ctx.Locals("partnerId").(uuid.UUID)
	if partnerId != uuid.Nil {
		return res.Forbidden(ctx, res.AlreadyRegisteredAsPartner)
	}

	return ctx.Next()
}
