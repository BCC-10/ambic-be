package code

import (
	"ambic/internal/domain/env"
	"crypto/rand"
	"encoding/base64"
)

type CodeIf interface {
	GenerateToken() (string, error)
}

type Code struct {
	TokenLength int
}

func NewCode(env *env.Env) CodeIf {
	return &Code{
		TokenLength: env.TokenLength,
	}
}

func (c *Code) GenerateToken() (string, error) {
	bytes := make([]byte, c.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)

	if len(token) > c.TokenLength {
		token = token[:c.TokenLength]
	}

	return token, nil
}
