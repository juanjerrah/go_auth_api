package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
)

type TokenRepository interface {
	StoreToken(ctx context.Context, token string, authCtx *auth.AuthContext, expiration time.Duration) error
	GetToken(ctx context.Context, token string) (*auth.AuthContext, error)
	DeleteToken(ctx context.Context, token string) error
	InvalidateUserTokens(ctx context.Context, userID string) error
	TokenExists(ctx context.Context, token string) (bool, error)
}

type RedisTokenRepository struct {
	client *redis.Client
	prefix string
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &RedisTokenRepository{
		client: client,
		prefix: "token:",
	}
}

func (r *RedisTokenRepository) StoreToken(ctx context.Context, token string, authCtx *auth.AuthContext, expiration time.Duration) error {
	key := r.getKey(token)
	
	authCtxBytes, err := json.Marshal(authCtx)
	if err != nil {
		return fmt.Errorf("failed to marshal auth context: %w", err)
	}

	err = r.client.Set(ctx, key, authCtxBytes, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %w", err)
	}

	// Também armazenar relação user -> tokens para invalidação em massa
	userTokensKey := r.getUserTokensKey(authCtx.UserID)
	err = r.client.SAdd(ctx, userTokensKey, key).Err()
	if err != nil {
		return fmt.Errorf("failed to store user token relation: %w", err)
	}

	// Definir expiração para o conjunto de tokens do usuário
	r.client.Expire(ctx, userTokensKey, expiration+time.Hour*24)

	return nil
}

func (r *RedisTokenRepository) GetToken(ctx context.Context, token string) (*auth.AuthContext, error) {
	key := r.getKey(token)
	
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get token from redis: %w", err)
	}

	var authCtx auth.AuthContext
	err = json.Unmarshal(data, &authCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth context: %w", err)
	}

	return &authCtx, nil
}

func (r *RedisTokenRepository) DeleteToken(ctx context.Context, token string) error {
	key := r.getKey(token)
	
	// Primeiro obter o auth context para remover da lista de tokens do usuário
	authCtx, err := r.GetToken(ctx, token)
	if err != nil {
		return err
	}

	userTokensKey := r.getUserTokensKey(authCtx.UserID)
	r.client.SRem(ctx, userTokensKey, key)

	err = r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete token from redis: %w", err)
	}

	return nil
}

func (r *RedisTokenRepository) InvalidateUserTokens(ctx context.Context, userID string) error {
	userTokensKey := r.getUserTokensKey(userID)
	
	tokens, err := r.client.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	if len(tokens) <= 0 {
		return nil
	}
	err = r.client.Del(ctx, tokens...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user tokens: %w", err)
	}

	// Remover o conjunto de tokens do usuário
	err = r.client.Del(ctx, userTokensKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user tokens set: %w", err)
	}

	return nil
}

func (r *RedisTokenRepository) TokenExists(ctx context.Context, token string) (bool, error) {
	key := r.getKey(token)
	
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token existence: %w", err)
	}

	return exists > 0, nil
}

func (r *RedisTokenRepository) getKey(token string) string {
	return r.prefix + token
}

func (r *RedisTokenRepository) getUserTokensKey(userID string) string {
	return "user_tokens:" + userID
}