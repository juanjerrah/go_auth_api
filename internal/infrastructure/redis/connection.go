package redis

import (
	"context"
	"log"

	"github.com/juanjerrah/go_auth_api/internal/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.URI,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis successfully")
	return client, nil
}

func DisconnectRedis(ctx context.Context, client *redis.Client) error {
	return client.Close()
}