package entities

import "github.com/fraineri/plexo_backend/internal/core/types"

// OTP represents the one-time password entity.
type OTP struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Code      string `json:"code"`
	Purpose   string `json:"purpose"`
	ExpiresAt int64  `json:"expires_at"`
	UsedAt    *int64 `json:"used_at,omitempty"`
	types.Model
}
