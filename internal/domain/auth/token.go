package auth

import (
	"context"
	"time"

	"github.com/juanjerrah/go_auth_api/pkg/types"
)

// Interface do repositório de tokens (agora definida aqui)
type TokenRepository interface {
	StoreToken(ctx context.Context, token string, authCtx *types.AuthContext, expiration time.Duration) error
	GetToken(ctx context.Context, token string) (*types.AuthContext, error)
	DeleteToken(ctx context.Context, token string) error
	InvalidateUserTokens(ctx context.Context, userID string) error
	TokenExists(ctx context.Context, token string) (bool, error)
}