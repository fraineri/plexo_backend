package persistance

import (
	"context"
	"fmt"
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

func (r *userRoleRepo) FindByUserID(ctx context.Context, userID string) ([]*entities.UserRole, error) {
	executor := r.uow.GetExecutor()
	query := `
		SELECT user_id, role_id, assigned_at
		FROM user_roles
		WHERE user_id = $1
	`
	rows, err := executor.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}()

	var userRoles []*entities.UserRole
	for rows.Next() {
		var userRole entities.UserRole
		var assignedAt int64
		if err := rows.Scan(&userRole.UserID, &userRole.RoleID, &assignedAt); err != nil {
			return nil, err
		}
		userRole.AssignedAt = time.Unix(assignedAt, 0)
		userRoles = append(userRoles, &userRole)
	}

	return userRoles, nil
}
