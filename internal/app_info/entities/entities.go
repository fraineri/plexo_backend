package entities

import "github.com/fraineri/plexo_backend/internal/core/types"

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
	types.Model
}
