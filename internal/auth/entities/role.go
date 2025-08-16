package entities

import "github.com/fraineri/plexo_backend/internal/core/types"

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	types.Model
}
