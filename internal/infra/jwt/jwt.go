package jwt

import (
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTIf interface {
	GenerateToken(userId uuid.UUID, isVerified bool, partnerId uuid.UUID, isVerifiedPartner bool) (string, error)
	ValidateToken(token string) (uuid.UUID, bool, uuid.UUID, bool, error)
}
type JWT struct {
	secretKey   string
	expiredTime time.Duration
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
	Id                uuid.UUID
	IsVerified        bool
	PartnerId         uuid.UUID
	IsVerifiedPartner bool
	jwt.RegisteredClaims
}

func (j *JWT) GenerateToken(userId uuid.UUID, isVerified bool, partnerId uuid.UUID, isVerifiedPartner bool) (string, error) {
	claim := Claims{
		Id:                userId,
		IsVerified:        isVerified,
		PartnerId:         partnerId,
		IsVerifiedPartner: isVerifiedPartner,
		RegisteredClaims:  jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiredTime))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (j *JWT) ValidateToken(tokenString string) (uuid.UUID, bool, uuid.UUID, bool, error) {
	claim := new(Claims)

	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return uuid.Nil, false, uuid.Nil, false, err
	}

	if !token.Valid {
		return uuid.Nil, false, uuid.Nil, false, errors.New(res.InvalidToken)
	}

	userId := claim.Id
	isVerified := claim.IsVerified
	partnerId := claim.PartnerId
	isVerifiedPartner := claim.IsVerifiedPartner

	return userId, isVerified, partnerId, isVerifiedPartner, nil
}
