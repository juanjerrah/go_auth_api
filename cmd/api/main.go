package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/juanjerrah/go_auth_api/internal/config"
	"github.com/juanjerrah/go_auth_api/internal/delivery/http/handlers"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
	"github.com/juanjerrah/go_auth_api/internal/infrastructure/mongodb"
	"github.com/juanjerrah/go_auth_api/internal/infrastructure/redis"
	"github.com/juanjerrah/go_auth_api/internal/utils"
	"github.com/juanjerrah/go_auth_api/pkg/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongoClient, err := mongodb.ConnectMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Connect to Redis
	redisClient, err := redis.ConnectRedis(&cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	// Initialize Infrastructure
	mongoDB := mongodb.NewMongoDB(mongoClient, cfg.MongoDB.Database)
	userRepo := mongodb.NewUserRepository(mongoDB.Database)
	tokenRepo := redis.NewTokenRepository(redisClient)

	// Initialize utilities
	passwordHasher := utils.NewBcryptPasswordHasher(12)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.TokenExpiresIn)
	mongoUtils := utils.NewMongoUtils()

	// Initialize Services
	userService := user.NewService(userRepo, passwordHasher, mongoUtils)
	authService := auth.NewAuthService(tokenRepo)

	// Initialize Gin
	router := gin.Default()

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	userHandler := handlers.NewUserHandler(userService)

	// Routes
	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)

		// Protected routes
		protected := api.Group("", middleware.AuthMiddleware(jwtManager, authService))
		{
			// User routes
			protected.PUT("/users/:id/password", userHandler.ChangePassword)
		}
	}

	// Start server
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}

}
