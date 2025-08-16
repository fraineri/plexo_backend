package entities

import "time"

type RolePermission struct {
	RoleID       string    `json:"role_id"`
	PermissionID string    `json:"permission_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
