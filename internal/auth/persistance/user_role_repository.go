package persistance

import (
	"context"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type userRoleRepo struct {
	uow uow.UnitOfWork
}

func NewUserRoleRepository(uow uow.UnitOfWork) UserRoleRepository {
	return &userRoleRepo{uow: uow}
}

func (r *userRoleRepo) Create(ctx context.Context, userRole *entities.UserRole) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		RETURNING assigned_at
	`
	var assignedAt int64
	err := executor.QueryRowContext(ctx, query,
		userRole.UserID,
		userRole.RoleID,
	).Scan(&assignedAt)
	if err != nil {
		return err
	}
	userRole.AssignedAt = time.Unix(assignedAt, 0)
	return nil
}
