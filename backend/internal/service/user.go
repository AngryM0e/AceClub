package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
)

type UserServiceInterface interface {
	RegisterUser(ctx context.Context, email, password string) (*domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// UserService manages user registration and authentication business logic
type UserService struct {
	repo   domain.UserRepository
	hasher PasswordHasher
}

func NewUserService(repo domain.UserRepository, hasher PasswordHasher) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *UserService) RegisterUser(ctx context.Context, email, password string) (*domain.User, error) {

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return nil, domain.NewValidationError("email", errors.New("invalid email format"))
	}

	if len(password) < 8 {
		return nil, domain.NewValidationError("password", errors.New("password length must be >= 8"))
	}

	if len(password) > 16 {
		return nil, domain.NewValidationError("password", errors.New("password length must be <= 16"))
	}

	existingUser, err := u.repo.GetByEmail(ctx, email)

	var ne *domain.NotFoundError
	if err != nil && !errors.As(err, &ne) {
		return nil, fmt.Errorf("check email uniqueness: %w", err)
	}

	if existingUser != nil {
		return nil, domain.NewConflictError("user", email)
	}

	hashedPass, err := u.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Email:    parsedEmail.Address,
		Password: hashedPass,
	}

	err = u.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
