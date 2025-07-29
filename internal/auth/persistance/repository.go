package persistance

import (
	"context"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error) // New method
	Update(ctx context.Context, user *entities.User) error
}

type OtpRepository interface {
	Create(ctx context.Context, otp *entities.OTP) error
	FindLatestByUserID(ctx context.Context, userID, purpose string) (*entities.OTP, error)
	CountRecentByUserID(ctx context.Context, userID, purpose string, since int64) (int, error)
	Update(ctx context.Context, otp *entities.OTP) error
	InvalidateAllUnused(ctx context.Context, userID, purpose string) error
}

type LoginAttemptRepository interface {
	Create(ctx context.Context, attempt *entities.LoginAttempt) error
	FindLastByUser(ctx context.Context, userID string, limit int) ([]*entities.LoginAttempt, error)
}
