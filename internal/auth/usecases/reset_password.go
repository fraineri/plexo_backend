package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type ResetPasswordInput struct {
	Email       string
	OTPCode     string
	NewPassword string
}

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, input ResetPasswordInput) error
}

type resetPassword struct {
	uow            uow.UnitOfWork
	userService    services.UserService
	authOtpService services.AuthOtpService
	userRepository persistance.UserRepository
}

func NewResetPassword(
	uow uow.UnitOfWork,
	userService services.UserService,
	authOtpService services.AuthOtpService,
	userRepository persistance.UserRepository,
) ResetPasswordUseCase {
	return &resetPassword{
		uow:            uow,
		userService:    userService,
		authOtpService: authOtpService,
		userRepository: userRepository,
	}
}

func (uc *resetPassword) Execute(ctx context.Context, input ResetPasswordInput) error {
	return uc.uow.Do(ctx, func() error {
		// 1. Find user
		user, err := uc.userRepository.FindByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find user by email: %w", err)
		}
		if user == nil {
			return fmt.Errorf("user with email %s not found", input.Email)
		}

		// 2. Verify the OTP is valid for this purpose.
		err = uc.authOtpService.VerifyOtp(ctx, user.ID, "PASSWORD_RESET", input.OTPCode)
		if err != nil {
			return fmt.Errorf("otp verification failed: %w", err)
		}

		// 3. Update the user's password.
		err = uc.userService.UpdatePassword(ctx, user.ID, input.NewPassword)
		if err != nil {
			return fmt.Errorf("password update failed: %w", err)
		}

		return nil
	})
}
