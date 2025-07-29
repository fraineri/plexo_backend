package services

import (
	"crypto/rand"
	"io"
)

type OtpService interface {
	Generate() (string, error)
}

type otpService struct{}

func NewOtpService() OtpService {
	return &otpService{}
}

func (s *otpService) Generate() (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	digits := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, digits, 6)
	if n != 6 {
		return "", err
	}
	for i := 0; i < len(digits); i++ {
		digits[i] = table[int(digits[i])%len(table)]
	}
	return string(digits), nil
}
