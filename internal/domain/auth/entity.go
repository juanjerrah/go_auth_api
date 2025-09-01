package auth

import "github.com/juanjerrah/go_auth_api/pkg/types"

// Re-export para manter compatibilidade
type Permission = types.Permission
type AuthContext = types.AuthContext

// Constantes re-exportadas
const (
	PermissionUserRead   = types.PermissionUserRead
	PermissionUserWrite  = types.PermissionUserWrite
	PermissionUserDelete = types.PermissionUserDelete
	PermissionAdminRead  = types.PermissionAdminRead
	PermissionAdminWrite = types.PermissionAdminWrite
)

// Funções de utilidade
func HasPermission(role types.Role, permission types.Permission) bool {
	return types.HasPermission(role, permission)
}