package persistance

import (
	"context"
	"database/sql"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type otpRepo struct {
	uow uow.UnitOfWork
}

func NewOtpRepository(uow uow.UnitOfWork) OtpRepository {
	return &otpRepo{uow: uow}
}

func (r *otpRepo) Create(ctx context.Context, otp *entities.OTP) error {
	executor := r.uow.GetExecutor()
	query := `
		INSERT INTO otps (user_id, code, purpose, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	var createdAt int64
	err := executor.QueryRowContext(ctx, query,
		otp.UserID,
		otp.Code,
		otp.Purpose,
		otp.ExpiresAt,
	).Scan(&otp.ID, &createdAt)
	if err != nil {
		return err
	}
	otp.CreatedAt = time.Unix(createdAt, 0)
	return nil
}

func (r *otpRepo) FindLatestByUserID(ctx context.Context, userID, purpose string) (*entities.OTP, error) {
	executor := r.uow.GetExecutor()
	query := `
		SELECT id, user_id, code, purpose, expires_at, used_at, created_at
		FROM otps
		WHERE user_id = $1 AND purpose = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	var otp entities.OTP
	var usedAt, createdAt sql.NullInt64
	err := executor.QueryRowContext(ctx, query, userID, purpose).Scan(
		&otp.ID,
		&otp.UserID,
		&otp.Code,
		&otp.Purpose,
		&otp.ExpiresAt,
		&usedAt,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if usedAt.Valid {
		otp.UsedAt = &usedAt.Int64
	}
	if createdAt.Valid {
		otp.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	return &otp, nil
}

func (r *otpRepo) CountRecentByUserID(ctx context.Context, userID, purpose string, since int64) (int, error) {
	executor := r.uow.GetExecutor()
	query := `
		SELECT COUNT(*)
		FROM otps
		WHERE user_id = $1 AND purpose = $2 AND created_at >= $3
	`
	var count int
	err := executor.QueryRowContext(ctx, query, userID, purpose, since).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *otpRepo) Update(ctx context.Context, otp *entities.OTP) error {
	executor := r.uow.GetExecutor()
	query := `UPDATE otps SET used_at = $1 WHERE id = $2`
	_, err := executor.ExecContext(ctx, query, otp.UsedAt, otp.ID)
	return err
}

func (r *otpRepo) InvalidateAllUnused(ctx context.Context, userID, purpose string) error {
	executor := r.uow.GetExecutor()
	query := `
		UPDATE otps
		SET expires_at = $1
		WHERE user_id = $2 AND purpose = $3 AND used_at IS NULL
	`
	now := time.Now().Unix()
	_, err := executor.ExecContext(ctx, query, now, userID, purpose)
	return err
}
