package persistance

import (
	"context"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type rolePermissionRepo struct {
	uow uow.UnitOfWork
}

func NewRolePermissionRepository(uow uow.UnitOfWork) RolePermissionRepository {
	return &rolePermissionRepo{uow: uow}
}

func (r *rolePermissionRepo) Create(ctx context.Context, rolePermission *entities.RolePermission) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO role_permissions (role_id, permission_id, status)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	var createdAt int64
	err := executor.QueryRowContext(ctx, query,
		rolePermission.RoleID,
		rolePermission.PermissionID,
		rolePermission.Status,
	).Scan(&createdAt)
	if err != nil {
		return err
	}
	rolePermission.CreatedAt = time.Unix(createdAt, 0)
	return nil
}
