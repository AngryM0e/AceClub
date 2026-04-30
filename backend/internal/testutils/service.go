package testutils

import (
	"context"
	"testing"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *MockHasher) Compare(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func SetupUserService(t *testing.T) (context.Context, *MockUserRepository, *MockHasher) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockHasher)

	t.Cleanup(func() {
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	return ctx, mockRepo, mockHasher
}
