package persistance

import (
	"context"
	"database/sql"
	"time"

	"github.com/fraineri/plexo_backend/internal/app_info/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type AppInfo = entities.AppInfo

type Repository interface {
	GetAppInfo(ctx context.Context) AppInfo
	Update(ctx context.Context, appInfo AppInfo) error
}

type AppInfoRepository struct {
	uow uow.UnitOfWork
}

func NewAppInfoRepository(uow uow.UnitOfWork) Repository {
	return &AppInfoRepository{uow: uow}
}

func (r *AppInfoRepository) GetAppInfo(ctx context.Context) AppInfo {
	executor := r.uow.GetExecutor()

	query := `SELECT name, version, status, created_at, updated_at, deleted_at FROM app_info LIMIT 1`
	row := executor.QueryRowContext(ctx, query)

	var appInfo AppInfo
	var createdAt, updatedAt, deletedAt sql.NullInt64
	err := row.Scan(
		&appInfo.Name,
		&appInfo.Version,
		&appInfo.Status,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return AppInfo{}
	}
	if createdAt.Valid {
		appInfo.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		appInfo.UpdatedAt = time.Unix(updatedAt.Int64, 0)
	}
	if deletedAt.Valid {
		appInfo.DeletedAt = time.Unix(deletedAt.Int64, 0)
	}
	return appInfo
}

func (r *AppInfoRepository) Update(ctx context.Context, appInfo AppInfo) error {
	executor := r.uow.GetExecutor()

	query := `UPDATE app_info SET version = $1, status = $2 WHERE name = $3`
	_, err := executor.ExecContext(ctx, query,
		appInfo.Version,
		appInfo.Status,
		appInfo.Name,
	)
	return err
}
