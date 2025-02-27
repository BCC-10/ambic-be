package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	authToken := ctx.GetReqHeaders()["Authorization"]

	if len(authToken) < 1 {
		return res.BadRequest(ctx, "token is not provided")
	}

	bearerToken := authToken[0]
	token := strings.Split(bearerToken, " ")

	userId, isVerified, err := m.jwt.ValidateToken(token[1])
	if err != nil {
		return res.Unauthorized(ctx, "invalid token")
	}

	ctx.Locals("userId", userId)
	ctx.Locals("isVerified", isVerified)

	return ctx.Next()
}
