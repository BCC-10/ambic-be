package code

import (
	"ambic/internal/domain/env"
	"crypto/rand"
	"fmt"
	"math/big"
)

type CodeIf interface {
	GenerateOTP() (string, error)
}

type Code struct {
	Length int
}

func NewCode(env *env.Env) CodeIf {
	return &Code{
		Length: env.OTPLength,
	}
}

func (c *Code) GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < c.Length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}
