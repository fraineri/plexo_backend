package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type ResendVerificationOtpInput struct {
	Email string
}

type ResendVerificationOtpUseCase interface {
	Execute(ctx context.Context, input ResendVerificationOtpInput) error
}

type resendVerificationOtp struct {
	uow                 uow.UnitOfWork
	userRepository      persistance.UserRepository
	authOtpService      services.AuthOtpService
	notificationService services.NotificationService
}

func NewResendVerificationOtp(
	uow uow.UnitOfWork,
	userRepository persistance.UserRepository,
	authOtpService services.AuthOtpService,
	notificationService services.NotificationService,
) ResendVerificationOtpUseCase {
	return &resendVerificationOtp{
		uow:                 uow,
		userRepository:      userRepository,
		authOtpService:      authOtpService,
		notificationService: notificationService,
	}
}

func (uc *resendVerificationOtp) Execute(ctx context.Context, input ResendVerificationOtpInput) error {
	return uc.uow.Do(ctx, func() error {
		// 1. Find user by email
		user, err := uc.userRepository.FindByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find user by email: %w", err)
		}
		if user == nil || user.Status == "ACTIVE" {
			// To prevent user enumeration, we don't return an error here.
			// The request will appear successful to the client, but no email will be sent.
			return nil
		}

		// 2. Resend the OTP
		otpCode, err := uc.authOtpService.ResendVerificationOtp(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to resend otp: %w", err)
		}

		// 3. Send notification
		if err := uc.notificationService.SendOTP(ctx, user.Email, otpCode); err != nil {
			return fmt.Errorf("failed to send otp email: %w", err)
		}

		return nil
	})
}
