package persistance

import (
	"context"
	"database/sql"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type userRepo struct {
	uow uow.UnitOfWork
}

func NewUserRepository(uow uow.UnitOfWork) UserRepository {
	return &userRepo{uow: uow}
}

func (r *userRepo) Create(ctx context.Context, user *entities.User) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO users (first_names, last_names, email, password_hash, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	var createdAt, updatedAt int64
	err := executor.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Status,
	).Scan(&user.ID, &createdAt, &updatedAt)
	if err != nil {
		return err
	}
	user.CreatedAt = time.Unix(createdAt, 0)
	user.UpdatedAt = time.Unix(updatedAt, 0)
	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	executor := r.uow.GetExecutor()
	query := `
		SELECT id, first_names, last_names, email, password_hash, status, blocked_until, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user entities.User
	var blockedUntil, createdAt, updatedAt sql.NullInt64
	err := executor.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Status,
		&blockedUntil,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if blockedUntil.Valid {
		user.BlockedUntil = &blockedUntil.Int64
	}
	if createdAt.Valid {
		user.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		user.UpdatedAt = time.Unix(updatedAt.Int64, 0)
	}
	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*entities.User, error) {
	executor := r.uow.GetExecutor()
	query := `
		SELECT id, first_names, last_names, email, password_hash, status, blocked_until, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entities.User
	var blockedUntil, createdAt, updatedAt sql.NullInt64
	err := executor.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Status,
		&blockedUntil,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	if blockedUntil.Valid {
		user.BlockedUntil = &blockedUntil.Int64
	}
	if createdAt.Valid {
		user.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		user.UpdatedAt = time.Unix(updatedAt.Int64, 0)
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *entities.User) error {
	executor := r.uow.GetExecutor()
	query := `
		UPDATE users
		SET status = $1, blocked_until = $2, password_hash = $3, updated_at = EXTRACT(EPOCH FROM NOW())
		WHERE id = $4
	`
	_, err := executor.ExecContext(ctx, query,
		user.Status,
		user.BlockedUntil,
		user.PasswordHash,
		user.ID,
	)
	return err
}
