package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) EnsureVerifiedPartner(ctx *fiber.Ctx) error {
	isVerifiedPartner := ctx.Locals("isVerifiedPartner").(bool)
	if !isVerifiedPartner {
		return res.Forbidden(ctx, res.PartnerNotVerified)
	}

	return ctx.Next()
}
