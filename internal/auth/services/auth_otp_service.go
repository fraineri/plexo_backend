package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
)

type AuthOtpService interface {
	CreateVerificationOtp(ctx context.Context, userID string) (string, error)
	VerifyOtp(ctx context.Context, userID, purpose, code string) error // New method
}

type authOtpService struct {
	otpRepository persistance.OtpRepository
}

func NewAuthOtpService(otpRepository persistance.OtpRepository) AuthOtpService {
	return &authOtpService{
		otpRepository: otpRepository,
	}
}

func (s *authOtpService) CreateVerificationOtp(ctx context.Context, userID string) (string, error) {
	otpCode, err := s.generateCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP code: %w", err)
	}

	otp := &entities.OTP{
		UserID:    userID,
		Code:      otpCode,
		Purpose:   "ACCOUNT_VERIFICATION",
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
	}

	if err := s.otpRepository.Create(ctx, otp); err != nil {
		return "", fmt.Errorf("failed to create OTP record: %w", err)
	}

	return otp.Code, nil
}

func (s *authOtpService) VerifyOtp(ctx context.Context, userID, purpose, code string) error {
	otp, err := s.otpRepository.FindLatestByUserID(ctx, userID, purpose)
	if err != nil {
		return fmt.Errorf("error finding otp - %w", err)
	}
	if otp == nil {
		return errors.New("invalid OTP - not found")
	}

	// Check 1: Code must match
	if otp.Code != code {
		return errors.New("invalid OTP - code mismatch")
	}

	// Check 2: OTP must not be used
	if otp.UsedAt != nil {
		return errors.New("invalid OTP - already used")
	}

	// Check 3: OTP must not be expired
	if time.Now().Unix() > otp.ExpiresAt {
		return errors.New("invalid OTP - expired")
	}

	// Mark OTP as used
	now := time.Now().Unix()
	otp.UsedAt = &now
	if err := s.otpRepository.Update(ctx, otp); err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	return nil
}

func (s *authOtpService) generateCode() (string, error) {
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
