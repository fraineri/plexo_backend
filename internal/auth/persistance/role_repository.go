package persistance

import (
	"context"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type roleRepo struct {
	uow uow.UnitOfWork
}

func NewRoleRepository(uow uow.UnitOfWork) RoleRepository {
	return &roleRepo{uow: uow}
}

func (r *roleRepo) Create(ctx context.Context, role *entities.Role) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO roles (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	var createdAt, updatedAt int64
	err := executor.QueryRowContext(ctx, query,
		role.Name,
		role.Description,
	).Scan(&role.ID, &createdAt, &updatedAt)
	if err != nil {
		return err
	}
	role.CreatedAt = time.Unix(createdAt, 0)
	role.UpdatedAt = time.Unix(updatedAt, 0)
	return nil
}
