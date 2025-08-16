package usecases

import (
	"context"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type AssignRoleToUserInput struct {
	UserID string
	RoleID string
}

type AssignRoleToUserUseCase interface {
	Execute(ctx context.Context, input AssignRoleToUserInput) error
}

type assignRoleToUser struct {
	uow                uow.UnitOfWork
	userRoleRepository persistance.UserRoleRepository
}

func NewAssignRoleToUser(
	uow uow.UnitOfWork,
	userRoleRepository persistance.UserRoleRepository,
) AssignRoleToUserUseCase {
	return &assignRoleToUser{
		uow:                uow,
		userRoleRepository: userRoleRepository,
	}
}

func (uc *assignRoleToUser) Execute(ctx context.Context, input AssignRoleToUserInput) error {
	return uc.uow.Do(ctx, func() error {
		userRole := &entities.UserRole{
			UserID: input.UserID,
			RoleID: input.RoleID,
		}
		return uc.userRoleRepository.Create(ctx, userRole)
	})
}
