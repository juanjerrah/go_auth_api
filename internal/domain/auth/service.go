package auth

import (
	"context"
	"errors"
	"time"

	"github.com/juanjerrah/go_auth_api/internal/domain/user"
	"github.com/juanjerrah/go_auth_api/pkg/types"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidRole      = errors.New("invalid role")
	ErrInvalidToken     = errors.New("invalid token")
)

type AuthService interface {
	HasPermission(role user.Role, permission types.Permission) bool
	GetUserPermissions(role user.Role) []types.Permission
	ValidateRole(role user.Role) error
	StoreToken(ctx context.Context, token string, authCtx *types.AuthContext, expiration time.Duration) error
	GetToken(ctx context.Context, token string) (*types.AuthContext, error)
	DeleteToken(ctx context.Context, token string) error
	InvalidateUserTokens(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (*types.AuthContext, error)
}

type authService struct {
	tokenRepo TokenRepository
}

func NewAuthService(tokenRepo TokenRepository) AuthService {
	return &authService{
		tokenRepo: tokenRepo,
	}
}

func (s *authService) HasPermission(role user.Role, permission types.Permission) bool {
	return types.HasPermission(role, permission)
}

func (s *authService) GetUserPermissions(role user.Role) []types.Permission {
	return types.RolePermissionMap[role]
}

func (s *authService) ValidateRole(role user.Role) error {
	if role != user.RoleUser && role != user.RoleAdmin {
		return ErrInvalidRole
	}
	return nil
}

func (s *authService) StoreToken(ctx context.Context, token string, authCtx *types.AuthContext, expiration time.Duration) error {
	return s.tokenRepo.StoreToken(ctx, token, authCtx, expiration)
}

func (s *authService) GetToken(ctx context.Context, token string) (*types.AuthContext, error) {
	return s.tokenRepo.GetToken(ctx, token)
}

func (s *authService) DeleteToken(ctx context.Context, token string) error {
	return s.tokenRepo.DeleteToken(ctx, token)
}

func (s *authService) InvalidateUserTokens(ctx context.Context, userID string) error {
	return s.tokenRepo.InvalidateUserTokens(ctx, userID)
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*types.AuthContext, error) {
	authCtx, err := s.tokenRepo.GetToken(ctx, token)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return authCtx, nil
}
