package usecases

import (
	"context"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type CreatePermissionInput struct {
	Action      string
	Description string
}

type CreatePermissionUseCase interface {
	Execute(ctx context.Context, input CreatePermissionInput) (*entities.Permission, error)
}

type createPermission struct {
	uow                  uow.UnitOfWork
	permissionRepository persistance.PermissionRepository
}

func NewCreatePermission(
	uow uow.UnitOfWork,
	permissionRepository persistance.PermissionRepository,
) CreatePermissionUseCase {
	return &createPermission{
		uow:                  uow,
		permissionRepository: permissionRepository,
	}
}

func (uc *createPermission) Execute(ctx context.Context, input CreatePermissionInput) (*entities.Permission, error) {
	var permission *entities.Permission
	err := uc.uow.Do(ctx, func() error {
		newPermission := &entities.Permission{
			Action:      input.Action,
			Description: input.Description,
		}
		if err := uc.permissionRepository.Create(ctx, newPermission); err != nil {
			return err
		}
		permission = newPermission
		return nil
	})
	if err != nil {
		return nil, err
	}
	return permission, nil
}
