package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type RequestPasswordResetInput struct {
	Email string
}

type RequestPasswordResetUseCase interface {
	Execute(ctx context.Context, input RequestPasswordResetInput) error
}

type requestPasswordReset struct {
	uow                 uow.UnitOfWork
	userRepository      persistance.UserRepository
	authOtpService      services.AuthOtpService
	notificationService services.NotificationService
}

func NewRequestPasswordReset(
	uow uow.UnitOfWork,
	userRepository persistance.UserRepository,
	authOtpService services.AuthOtpService,
	notificationService services.NotificationService,
) RequestPasswordResetUseCase {
	return &requestPasswordReset{
		uow:                 uow,
		userRepository:      userRepository,
		authOtpService:      authOtpService,
		notificationService: notificationService,
	}
}

func (uc *requestPasswordReset) Execute(ctx context.Context, input RequestPasswordResetInput) error {
	return uc.uow.Do(ctx, func() error {
		user, err := uc.userRepository.FindByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find user by email: %w", err)
		}

		// To prevent user enumeration, we do not return an error if the user is not found or not active.
		// The request appears successful, but no email is sent.
		if user == nil || user.Status != "ACTIVE" {
			return nil
		}

		// Create and send the password reset OTP.
		otpCode, err := uc.authOtpService.CreatePasswordResetOtp(ctx, user.ID)
		if err != nil {
			// This error (e.g., rate limit) is returned to the client, which is acceptable here.
			return err
		}

		if err := uc.notificationService.SendOTP(ctx, user.Email, otpCode); err != nil {
			return fmt.Errorf("failed to send otp email: %w", err)
		}

		return nil
	})
}
