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
	"github.com/fraineri/plexo_backend/internal/core/settings"
)

type AuthOtpService interface {
	CreateVerificationOtp(ctx context.Context, userID string) (string, error)
	ResendVerificationOtp(ctx context.Context, userID string) (string, error)
	CreatePasswordResetOtp(ctx context.Context, userID string) (string, error)
	VerifyOtp(ctx context.Context, userID, purpose, code string) error
}

type authOtpService struct {
	otpRepository persistance.OtpRepository
	settings      settings.OTPSettings
}

func NewAuthOtpService(otpRepository persistance.OtpRepository, settings settings.OTPSettings) AuthOtpService {
	return &authOtpService{
		otpRepository: otpRepository,
		settings:      settings,
	}
}

func (s *authOtpService) CreateVerificationOtp(ctx context.Context, userID string) (string, error) {
	purpose := "ACCOUNT_VERIFICATION"
	if err := s.checkRateLimit(ctx, userID, purpose); err != nil {
		return "", err
	}
	return s.generateAndSaveOtp(ctx, userID, purpose)
}

func (s *authOtpService) ResendVerificationOtp(ctx context.Context, userID string) (string, error) {
	purpose := "ACCOUNT_VERIFICATION"

	// Step 1: Check rate limit
	if err := s.checkRateLimit(ctx, userID, purpose); err != nil {
		return "", err
	}

	// Step 2: Invalidate all previous OTPs for this purpose
	if err := s.otpRepository.InvalidateAllUnused(ctx, userID, purpose); err != nil {
		return "", fmt.Errorf("failed to invalidate old otps: %w", err)
	}

	// Step 3: Generate and save a new OTP
	return s.generateAndSaveOtp(ctx, userID, purpose)
}

func (s *authOtpService) CreatePasswordResetOtp(ctx context.Context, userID string) (string, error) {
	purpose := "PASSWORD_RESET"

	// Step 1: Check rate limit for this purpose
	if err := s.checkRateLimit(ctx, userID, purpose); err != nil {
		return "", err
	}

	// Step 2: Invalidate all previous OTPs for this purpose
	if err := s.otpRepository.InvalidateAllUnused(ctx, userID, purpose); err != nil {
		return "", fmt.Errorf("failed to invalidate old password reset otps: %w", err)
	}

	// Step 3: Generate and save a new OTP
	return s.generateAndSaveOtp(ctx, userID, purpose)
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

func (s *authOtpService) checkRateLimit(ctx context.Context, userID, purpose string) error {
	since := time.Now().Add(-time.Duration(s.settings.RateLimitMinutes) * time.Minute).Unix()
	count, err := s.otpRepository.CountRecentByUserID(ctx, userID, purpose, since)
	if err != nil {
		return fmt.Errorf("failed to count recent otps: %w", err)
	}

	if count >= s.settings.RateLimitCount {
		return ErrRateLimitExceeded
	}
	return nil
}

func (s *authOtpService) generateAndSaveOtp(ctx context.Context, userID, purpose string) (string, error) {
	otpCode, err := s.generateCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP code: %w", err)
	}

	otp := &entities.OTP{
		UserID:    userID,
		Code:      otpCode,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
	}

	if err := s.otpRepository.Create(ctx, otp); err != nil {
		return "", fmt.Errorf("failed to create OTP record: %w", err)
	}

	return otp.Code, nil
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
