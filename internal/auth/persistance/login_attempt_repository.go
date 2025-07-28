package persistance

import (
	"context"
	"fmt"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type loginAttemptRepo struct {
	uow uow.UnitOfWork
}

func NewLoginAttemptRepository(uow uow.UnitOfWork) LoginAttemptRepository {
	return &loginAttemptRepo{uow: uow}
}

func (r *loginAttemptRepo) Create(ctx context.Context, attempt *entities.LoginAttempt) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO login_attempts (user_id, is_successful, ip_address, user_agent)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	var createdAt int64
	err := executor.QueryRowContext(ctx, query,
		attempt.UserID,
		attempt.IsSuccessful,
		attempt.IPAddress,
		attempt.UserAgent,
	).Scan(&attempt.ID, &createdAt)
	if err != nil {
		return err
	}
	attempt.CreatedAt = time.Unix(createdAt, 0)
	return nil
}

func (r *loginAttemptRepo) FindLastByUser(ctx context.Context, userID string, limit int) ([]*entities.LoginAttempt, error) {
	executor := r.uow.GetExecutor()
	query := `
		SELECT id, user_id, is_successful, ip_address, user_agent, created_at
		FROM login_attempts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := executor.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}()

	var attempts []*entities.LoginAttempt
	for rows.Next() {
		var attempt entities.LoginAttempt
		var createdAt int64
		if err := rows.Scan(&attempt.ID, &attempt.UserID, &attempt.IsSuccessful, &attempt.IPAddress, &attempt.UserAgent, &createdAt); err != nil {
			return nil, err
		}
		attempt.CreatedAt = time.Unix(createdAt, 0)
		attempts = append(attempts, &attempt)
	}

	return attempts, nil
}
