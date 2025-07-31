package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type VerifyAccountInput struct {
	Email   string
	OTPCode string
}

type VerifyAccountUseCase interface {
	Execute(ctx context.Context, input VerifyAccountInput) error
}

type verifyAccount struct {
	uow            uow.UnitOfWork
	userService    services.UserService
	authOtpService services.AuthOtpService
	userRepository persistance.UserRepository
}

func NewVerifyAccount(
	uow uow.UnitOfWork,
	userService services.UserService,
	authOtpService services.AuthOtpService,
	userRepository persistance.UserRepository,
) VerifyAccountUseCase {
	return &verifyAccount{
		uow:            uow,
		userService:    userService,
		authOtpService: authOtpService,
		userRepository: userRepository,
	}
}

func (uc *verifyAccount) Execute(ctx context.Context, input VerifyAccountInput) error {
	err := uc.uow.Do(ctx, func() error {
		// 1. Find the user by email to get their ID
		user, err := uc.userRepository.FindByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find user by email: %w", err)
		}
		if user == nil {
			return fmt.Errorf("user with email %s not found", input.Email)
		}

		// 2. Verify the OTP is valid via the AuthOtpService
		err = uc.authOtpService.VerifyOtp(ctx, user.ID, "ACCOUNT_VERIFICATION", input.OTPCode)
		if err != nil {
			// This will return specific errors like "expired", "already used", etc.
			return fmt.Errorf("otp verification failed: %w", err)
		}

		// 3. Activate the user via the UserService
		err = uc.userService.ActivateUser(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("account activation failed: %w", err)
		}

		return nil
	})

	return err
}
