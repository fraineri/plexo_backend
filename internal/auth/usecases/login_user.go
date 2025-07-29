package usecases

import (
	"context"

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
	var token string

	err := uc.uow.Do(ctx, func() error {
		// 1. Authenticate user
		user, err := uc.userService.Login(ctx, input.Email, input.Password)
		if err != nil {
			return err // Returns specific errors like "invalid credentials" or "account blocked"
		}

		// 2. Generate JWT token
		generatedToken, err := uc.jwtService.GenerateToken(user)
		if err != nil {
			return err
		}
		token = generatedToken
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{Token: token}, nil
}
