package redis

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/juanjerrah/go_auth_api/internal/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	options := &redis.Options{
		Addr:     cfg.URI,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	// Configure SSL/TLS if enabled
	if cfg.UseSSL {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(options)

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
