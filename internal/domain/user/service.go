package user

import (
	"context"
	"errors"
	"time"

	"github.com/juanjerrah/go_auth_api/pkg/common"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

type Service interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) error
	DeleteUser(ctx context.Context, id string) error
	Authenticate(ctx context.Context, email, password string) (*User, error)
	ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error
}

type service struct {
	repo       Repository
	hasher     common.PasswordHasher
	mongoUtils common.MongoUtils
}

func NewService(repo Repository, hasher common.PasswordHasher, mongoUtils common.MongoUtils) Service {
	return &service{
		repo:       repo,
		hasher:     hasher,
		mongoUtils: mongoUtils,
	}
}

// ChangePassword implements Service.
func (s *service) ChangePassword(ctx context.Context, id string, oldPassword string, newPassword string) error {

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.hasher.Verify(oldPassword, user.Password); err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// Authenticate implements Service.
func (s *service) Authenticate(ctx context.Context, email string, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	if err := s.hasher.Verify(password, user.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// CreateUser implements Service.
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	exist, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrEmailAlreadyInUse
	}

	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	var user = &User{
		ID:        s.mongoUtils.GenerateObjectID(),
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.toResponse(user), nil
}

// DeleteUser implements Service.
func (s *service) DeleteUser(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetUserByEmail implements Service.
func (s *service) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return s.toResponse(user), nil
}

// GetUserByID implements Service.
func (s *service) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return s.toResponse(user), nil
}

// UpdateUser implements Service.
func (s *service) UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *service) toResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
