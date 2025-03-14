package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	authToken := ctx.GetReqHeaders()["Authorization"]

	if len(authToken) < 1 {
		return res.Unauthorized(ctx, res.MissingToken)
	}

	bearerToken := authToken[0]

	token := strings.Split(bearerToken, " ")

	if len(token) < 2 {
		return res.Unauthorized(ctx, res.InvalidTokenFormat)
	}

	userId, isVerified, partnerId, isVerifiedPartner, err := m.jwt.ValidateToken(token[1])
	if err != nil {
		return res.Unauthorized(ctx, res.InvalidToken)
	}

	if !isVerified {
		return res.Unauthorized(ctx, res.UserNotVerified)
	}

	ctx.Locals("userId", userId)
	ctx.Locals("partnerId", partnerId)
	ctx.Locals("isVerifiedPartner", isVerifiedPartner)

	return ctx.Next()
}
