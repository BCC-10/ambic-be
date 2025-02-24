package jwt

import (
	"ambic/internal/domain/env"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTIf interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}
type JWT struct {
	secretKey   string
	expiredTime int
}

func NewJwt(env *env.Env) JWTIf {
	secretKey := env.JWTSecret
	expiresTime := env.JWTExpires

	return &JWT{
		secretKey:   secretKey,
		expiredTime: expiresTime,
	}
}

type Claims struct {
	Id uuid.UUID
	jwt.RegisteredClaims
}

func (j *JWT) GenerateToken(userId uuid.UUID) (string, error) {
	claim := Claims{
		Id:               userId,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expiredTime)))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (j *JWT) ValidateToken(tokenString string) (uuid.UUID, error) {
	claim := new(Claims)

	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	userId := claim.Id

	return userId, nil
}
