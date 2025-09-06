// @title           Go Auth API
// @version         1.0
// @description     API de autenticação e autorização com Go
// @termsOfService  http://swagger.io/terms/

// @contact.name   Juan Jerrah
// @contact.url    https://github.com/juanjerrah
// @contact.email  seu-email@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/juanjerrah/go_auth_api/docs" // swagger docs gerado automaticamente
	"github.com/juanjerrah/go_auth_api/internal/config"
	"github.com/juanjerrah/go_auth_api/internal/delivery/http/handlers"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
	"github.com/juanjerrah/go_auth_api/internal/infrastructure/mongodb"
	"github.com/juanjerrah/go_auth_api/internal/infrastructure/redis"
	"github.com/juanjerrah/go_auth_api/internal/utils"
	"github.com/juanjerrah/go_auth_api/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: .env file not found or error loading it. Using default values or system environment variables.")
	} else {
		log.Println("Environment variables loaded from .env file")
	}
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
	authHandler := handlers.NewAuthHandler(userService, jwtManager, authService)
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

	// Configurando o Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}

}
