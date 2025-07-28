package entities

import "github.com/fraineri/plexo_backend/internal/core/types"

// LoginAttempt represents a record of a login attempt.
type LoginAttempt struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	IsSuccessful bool   `json:"is_successful"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	types.Model
}
