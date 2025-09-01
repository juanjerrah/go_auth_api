package types

import (

	"slices"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
)

type Role = user.Role

type Permission string

const (
	PermissionUserRead   Permission = "user:read"
	PermissionUserWrite  Permission = "user:write"
	PermissionUserDelete Permission = "user:delete"
	PermissionAdminRead  Permission = "admin:read"
	PermissionAdminWrite Permission = "admin:write"
)

type AuthContext struct {
	UserID      string
	Email       string
	Role        Role
	Permissions []Permission
}

// Mapa de permissões por role (agora em types para evitar cycle)
var RolePermissionMap = map[Role][]Permission{
	user.RoleUser: {
		PermissionUserRead,
	},
	user.RoleAdmin: {
		PermissionUserRead,
		PermissionUserWrite,
		PermissionUserDelete,
		PermissionAdminRead,
		PermissionAdminWrite,
	},
}

func HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissionMap[role]
	if !exists {
		return false
	}

	return slices.Contains(permissions, permission)
}