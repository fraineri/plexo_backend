package usecases

import (
	"context"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type AssignPermissionToRoleInput struct {
	RoleID       string
	PermissionID string
	Status       string
}

type AssignPermissionToRoleUseCase interface {
	Execute(ctx context.Context, input AssignPermissionToRoleInput) error
}

type assignPermissionToRole struct {
	uow                      uow.UnitOfWork
	rolePermissionRepository persistance.RolePermissionRepository
}

func NewAssignPermissionToRole(
	uow uow.UnitOfWork,
	rolePermissionRepository persistance.RolePermissionRepository,
) AssignPermissionToRoleUseCase {
	return &assignPermissionToRole{
		uow:                      uow,
		rolePermissionRepository: rolePermissionRepository,
	}
}

func (uc *assignPermissionToRole) Execute(ctx context.Context, input AssignPermissionToRoleInput) error {
	return uc.uow.Do(ctx, func() error {
		rolePermission := &entities.RolePermission{
			RoleID:       input.RoleID,
			PermissionID: input.PermissionID,
			Status:       input.Status,
		}
		return uc.rolePermissionRepository.Create(ctx, rolePermission)
	})
}
