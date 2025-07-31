package usecases

import (
	"context"
	"errors"

	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	Token string `json:"token"`
}

type LoginUserUseCase interface {
	Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error)
}

type loginUser struct {
	uow         uow.UnitOfWork
	userService services.UserService
	jwtService  services.JWTService
}

func NewLoginUser(
	uow uow.UnitOfWork,
	userService services.UserService,
	jwtService services.JWTService,
) LoginUserUseCase {
	return &loginUser{
		uow:         uow,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (uc *loginUser) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error) {
	// 1. Authenticate user credentials. This happens outside any transaction.
	user, authErr := uc.userService.Authenticate(ctx, input.Email, input.Password)

	// 2. Handle failed authentication.
	if authErr != nil {
		if errors.Is(authErr, services.ErrInvalidCredentials) && user != nil {
			var recordFailureBusinessOutcome error

			// This transaction's purpose is to record the failure and MUST commit.
			txSystemErr := uc.uow.Do(ctx, func() error {
				recordFailureBusinessOutcome = uc.userService.RecordFailedLogin(ctx, user.ID)

				// Check if the error is NOT a domain error (e.g., DB connection issue).
				if recordFailureBusinessOutcome != nil && !errors.Is(recordFailureBusinessOutcome, services.ErrAccountLockout) {
					return recordFailureBusinessOutcome
				}

				// We return nil to ensure the UnitOfWork COMMITS the transaction.
				return nil
			})

			if txSystemErr != nil {
				return nil, txSystemErr
			}

			// After the transaction has committed, we check the business outcome.
			if errors.Is(recordFailureBusinessOutcome, services.ErrAccountLockout) {
				return nil, services.ErrAccountLockout
			}
		}

		return nil, authErr
	}

	// 3. Handle successful authentication.
	var token string
	txErr := uc.uow.Do(ctx, func() error {
		if err := uc.userService.RecordSuccessfulLogin(ctx, user.ID); err != nil {
			return err
		}
		generatedToken, err := uc.jwtService.GenerateToken(user)
		if err != nil {
			return err
		}
		token = generatedToken
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &LoginUserOutput{Token: token}, nil
}
