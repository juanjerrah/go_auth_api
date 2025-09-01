package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
)

func AuthMiddleware(jwtManager *auth.JWTManager, authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Verificar assinatura do token
		if _, err := jwtManager.VerifyToken(tokenString); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
			c.Abort()
			return
		}

		// Validar token no Redis
		authCtx, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Adicionar informações de autenticação ao contexto
		c.Set("authContext", authCtx)

		c.Next()
	}
}

func RoleBasedAuthMiddleware(requiredRole user.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := c.Get("authContext")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		authCtx := authContext.(*auth.AuthContext)
		if authCtx.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func PermissionMiddleware(requiredPermission auth.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := c.Get("authContext")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		authCtx := authContext.(*auth.AuthContext)
		if !auth.HasPermission(authCtx.Role, requiredPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
