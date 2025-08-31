package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort     string
	JWTSecret      string
	TokenExpiresIn time.Duration
	MongoDB        MongoDBConfig
	Redis          RedisConfig
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type RedisConfig struct {
	URI      string
	Password string
	DB       int
	Timeout  time.Duration
}

func LoadConfig() *Config {
	tokenExpiresIn, _ := strconv.Atoi(getEnv("TOKEN_EXPIRES_IN", "3600"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiresIn: time.Duration(tokenExpiresIn) * time.Second,
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "webapi"),
			Timeout:  10 * time.Second,
		},
		Redis: RedisConfig{
			URI:      getEnv("REDIS_URI", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
			Timeout:  5 * time.Second,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Using default value for %s: %s", key, defaultValue)
		return defaultValue
	}
	return value
}