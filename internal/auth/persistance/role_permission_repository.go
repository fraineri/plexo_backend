package persistance

import (
	"context"
	"fmt"
	"strings"
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

func (r *rolePermissionRepo) FindPermissionsByRoleIDs(ctx context.Context, roleIDs []string) ([]*entities.Permission, error) {
	executor := r.uow.GetExecutor()

	// Create a placeholder string for the IN clause
	placeholders := make([]string, len(roleIDs))
	args := make([]interface{}, len(roleIDs))
	for i, id := range roleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`
		SELECT p.id, p.action, p.description
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}()

	var permissions []*entities.Permission
	for rows.Next() {
		var permission entities.Permission
		if err := rows.Scan(&permission.ID, &permission.Action, &permission.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}
