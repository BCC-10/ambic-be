package code

import (
	"ambic/internal/domain/env"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

type CodeIf interface {
	GenerateOTP() (string, error)
	GenerateToken() (string, error)
}

type Code struct {
	OTPLength   int
	TokenLength int
}

func NewCode(env *env.Env) CodeIf {
	return &Code{
		OTPLength:   env.OTPLength,
		TokenLength: env.TokenLength,
	}
}

func (c *Code) GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < c.OTPLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
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
