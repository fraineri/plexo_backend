package persistance

import (
	"context"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type permissionRepo struct {
	uow uow.UnitOfWork
}

func NewPermissionRepository(uow uow.UnitOfWork) PermissionRepository {
	return &permissionRepo{uow: uow}
}

func (r *permissionRepo) Create(ctx context.Context, permission *entities.Permission) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO permissions (action, description)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	var createdAt int64
	err := executor.QueryRowContext(ctx, query,
		permission.Action,
		permission.Description,
	).Scan(&permission.ID, &createdAt)
	if err != nil {
		return err
	}
	permission.CreatedAt = time.Unix(createdAt, 0)
	return nil
}
