package user

import (
	"context"
	"errors"
	"time"

	"github.com/juanjerrah/go_auth_api/internal/domain/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

type service struct {
	repo   Repository
	hasher hash.PasswordHasher
}

func NewService(repo Repository, hasher hash.PasswordHasher) Service {
	return &service{
		repo:   repo,
		hasher: hasher,
	}
}

// Authenticate implements Service.
func (s *service) Authenticate(ctx context.Context, email string, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	if err := s.hasher.Compare(user.Password, password); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// CreateUser implements Service.
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	exist, err := s.repo.ExistsByEmail(ctx, req.Email);
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
		ID:       primitive.NewObjectID(),
		Email:    req.Email,
		Password: hashedPassword,
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
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, userID); err != nil {
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
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return s.toResponse(user), nil
}

// UpdateUser implements Service.
func (s *service) UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrUserNotFound
	}

	user, err := s.repo.FindByID(ctx, userID)
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
