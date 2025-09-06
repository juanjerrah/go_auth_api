package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juanjerrah/go_auth_api/internal/delivery/http/handlers"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
	"github.com/juanjerrah/go_auth_api/pkg/middleware"
)

func SetupRoutes(router *gin.Engine, userService user.Service, authService auth.AuthService, jwtManager *auth.JWTManager) {
	// Handlers
	authHandler := handlers.NewAuthHandler(userService, jwtManager, authService)
	userHandler := handlers.NewUserHandler(userService)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtManager, authService))
	{
		// User routes
		userRoutes := protected.Group("/users")
		{
			userRoutes.GET("/profile", userHandler.GetUserProfile)
			userRoutes.PUT("/:id", userHandler.UpdateUser)
			userRoutes.DELETE("/:id", userHandler.DeleteUser)
		}

		// Admin only routes
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(middleware.PermissionMiddleware(auth.PermissionAdminRead))
		{
			//adminRoutes.GET("/users", userHandler.ListUsers)
			adminRoutes.GET("/users/:id", userHandler.GetUserByID)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}