package middleware

import (
	"ambic/internal/infra/jwt"
)

type MiddlewareIf interface {
}

type Middleware struct {
	jwt jwt.JWTIf
}

func NewMiddleware(jwt jwt.JWTIf) MiddlewareIf {
	return &Middleware{
		jwt: jwt,
	}
}
