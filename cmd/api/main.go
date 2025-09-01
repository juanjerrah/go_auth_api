package api

import (
	"context"
	"log"

	"github.com/juanjerrah/go_auth_api/internal/config"
	"github.com/juanjerrah/go_auth_api/internal/infrastructure/mongodb"
	"github.com/juanjerrah/go_auth_api/internal/utils"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongoClient, err := mongodb.ConnectMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatal(err)
	defer mongoClient.Disconnect(nil)
	
	// Connect to Redis
	redisClient, err := redis.ConnectRedis(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()
	
	// Initialize Infrastructure
	mongoDB := mongodb.NewMongoDB(mongoClient, cfg.MongoDB.Database)
	userRepo := mongodb.NewUserRepository(mongoDB)

	// Initialize utilities
	passwordHasher := utils.NewBcryptPasswordHasher(12)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.TokenExpiresIn)
	mongoUtils := utils.NewMongoUtils()

	// Initialize Services
	userService := user.NewService(userRepo, passwordHasher, mongoUtils)

	// Initialize Gin
	router := gin.Default()

}