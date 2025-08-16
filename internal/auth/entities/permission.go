package entities

import "github.com/fraineri/plexo_backend/internal/core/types"

type Permission struct {
	ID          string `json:"id"`
	Action      string `json:"action"`
	Description string `json:"description"`
	types.Model
}
