package entities

import "github.com/fraineri/plexo_backend/internal/core/types"

// User represents the user entity in the database.
type User struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_names"`
	LastName     string `json:"last_names"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // Omitted from JSON responses
	Status       string `json:"status"`
	BlockedUntil *int64 `json:"blocked_until,omitempty"`
	types.Model
}
