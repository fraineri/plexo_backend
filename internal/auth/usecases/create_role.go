package usecases

import (
	"context"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type CreateRoleInput struct {
	Name        string
	Description string
}

type CreateRoleUseCase interface {
	Execute(ctx context.Context, input CreateRoleInput) (*entities.Role, error)
}

type createRole struct {
	uow            uow.UnitOfWork
	roleRepository persistance.RoleRepository
}

func NewCreateRole(
	uow uow.UnitOfWork,
	roleRepository persistance.RoleRepository,
) CreateRoleUseCase {
	return &createRole{
		uow:            uow,
		roleRepository: roleRepository,
	}
}

func (uc *createRole) Execute(ctx context.Context, input CreateRoleInput) (*entities.Role, error) {
	var role *entities.Role
	err := uc.uow.Do(ctx, func() error {
		newRole := &entities.Role{
			Name:        input.Name,
			Description: input.Description,
		}
		if err := uc.roleRepository.Create(ctx, newRole); err != nil {
			return err
		}
		role = newRole
		return nil
	})
	if err != nil {
		return nil, err
	}
	return role, nil
}
