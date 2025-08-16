package services

import (
	"context"
	"strings"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
)

type AuthorizationService interface {
	Can(ctx context.Context, userID string, requiredPermission string) (bool, error)
}

type authorizationService struct {
	userRoleRepo       persistance.UserRoleRepository
	rolePermissionRepo persistance.RolePermissionRepository
}

func NewAuthorizationService(
	userRoleRepo persistance.UserRoleRepository,
	rolePermissionRepo persistance.RolePermissionRepository,
) AuthorizationService {
	return &authorizationService{
		userRoleRepo:       userRoleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

func (s *authorizationService) Can(ctx context.Context, userID string, requiredPermission string) (bool, error) {
	// 1. Get user roles
	userRoles, err := s.userRoleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	if len(userRoles) == 0 {
		return false, nil // No roles, no access
	}

	var roleIDs []string
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// 2. Get all permissions for those roles
	permissions, err := s.rolePermissionRepo.FindPermissionsByRoleIDs(ctx, roleIDs)
	if err != nil {
		return false, err
	}

	// 3. Evaluate permissions
	userPermissions := make(map[string]string)
	for _, p := range permissions {
		// Apply Deny-First rule
		if val, ok := userPermissions[p.Action]; ok && val == "DENY" {
			continue
		}
		userPermissions[p.Action] = "ALLOW" // Assuming status is implicitly ALLOW for now
	}

	// Check for exact match
	if status, ok := userPermissions[requiredPermission]; ok && status == "ALLOW" {
		return true, nil
	}

	// Check for wildcard match
	for permission, status := range userPermissions {
		if strings.HasSuffix(permission, ".*") {
			prefix := strings.TrimSuffix(permission, ".*")
			if strings.HasPrefix(requiredPermission, prefix) && status == "ALLOW" {
				return true, nil
			}
		}
	}

	return false, nil
}
