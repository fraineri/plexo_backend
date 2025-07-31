package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type RegisterUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type RegisterUserUseCase interface {
	Execute(ctx context.Context, input RegisterUserInput) (*entities.User, error)
}

type registerUser struct {
	uow                 uow.UnitOfWork
	userService         services.UserService
	authOtpService      services.AuthOtpService
	hashingService      services.HashingService
	notificationService services.NotificationService
}

func NewRegisterUser(
	uow uow.UnitOfWork,
	userService services.UserService,
	authOtpService services.AuthOtpService,
	hashingService services.HashingService,
	notificationService services.NotificationService,
) RegisterUserUseCase {
	return &registerUser{
		uow:                 uow,
		userService:         userService,
		authOtpService:      authOtpService,
		hashingService:      hashingService,
		notificationService: notificationService,
	}
}

func (uc *registerUser) Execute(ctx context.Context, input RegisterUserInput) (*entities.User, error) {
	var user *entities.User

	err := uc.uow.Do(ctx, func() error {
		// 1. Check if email exists via UserService
		exists, err := uc.userService.CheckIfEmailExists(ctx, input.Email)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("email '%s' is already in use", input.Email)
		}

		// 2. Hash the password
		hashedPassword, err := uc.hashingService.HashPassword(input.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// 3. Create the user record via UserService
		newUser, err := uc.userService.CreateUser(ctx, input.FirstName, input.LastName, input.Email, hashedPassword)
		if err != nil {
			return err
		}

		// 4. Generate and store OTP via AuthOtpService
		otpCode, err := uc.authOtpService.CreateVerificationOtp(ctx, newUser.ID)
		if err != nil {
			return err
		}

		// 5. Send OTP to user's email
		if err := uc.notificationService.SendOTP(ctx, newUser.Email, otpCode); err != nil {
			return fmt.Errorf("failed to send OTP email: %w", err)
		}

		user = newUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
