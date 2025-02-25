package code

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type CodeIf interface {
	GenerateOTP(length int) (string, error)
}

type Code struct{}

func NewCode() CodeIf {
	return &Code{}
}

func (c *Code) GenerateOTP(length int) (string, error) {
	otp := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}
