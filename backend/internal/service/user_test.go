package service

import (
	"errors"
	"testing"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/AngryM0e/AceClub/Backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser_Success(t *testing.T) {
	ctx, mockRepo, mockHasher := testutils.SetupUserService(t)
	svc := NewUserService(mockRepo, mockHasher)

	email := "henry@email.com"
	password := "password123"
	hashedPassword := "hashed_password_123"

	mockRepo.On("GetByEmail", ctx, email).Return(nil, nil)

	mockHasher.On("Hash", password).Return(hashedPassword, nil)

	mockRepo.On("Create", ctx, mock.MatchedBy(func(user *domain.User) bool {
		return user.Email == email && user.Password == hashedPassword
	})).Return(nil)

	user, err := svc.RegisterUser(ctx, email, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	ctx, mockRepo, mockHasher := testutils.SetupUserService(t)
	svc := NewUserService(mockRepo, mockHasher)

	email := "henry@email.com"
	password := "password123"

	existingUser := &domain.User{
		ID:       123,
		Email:    email,
		Password: "existing_hashed_password",
	}

	mockRepo.On("GetByEmail", ctx, email).Return(existingUser, nil)

	user, err := svc.RegisterUser(ctx, email, password)

	assert.Nil(t, user)
	assert.Error(t, err)

	var ce *domain.ConflictError
	assert.ErrorAs(t, err, &ce)

	assert.Contains(t, err.Error(), "user with henry@email.com already exists")
	assert.Contains(t, err.Error(), email)

	mockHasher.AssertNotCalled(t, "Hash", mock.Anything)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

	mockRepo.AssertCalled(t, "GetByEmail", ctx, email)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	ctx, mockRepo, mockHasher := testutils.SetupUserService(t)
	svc := NewUserService(mockRepo, mockHasher)

	tests := []struct {
		name          string
		email         string
		password      string
		expectedField string
		expectedMsg   string
	}{
		{
			name:          "empty email",
			email:         "",
			password:      "password123",
			expectedField: "email",
			expectedMsg:   "invalid email format",
		},
		{
			name:          "email without @ symbol",
			email:         "henryemail.com",
			password:      "password123",
			expectedField: "email",
			expectedMsg:   "invalid email format",
		},
		{
			name:          "email with spaces",
			email:         "test @example.com",
			password:      "password123",
			expectedField: "email",
			expectedMsg:   "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.RegisterUser(ctx, tt.email, tt.password)

			assert.Nil(t, user)
			assert.Error(t, err)

			var ve *domain.ValidationError
			if assert.ErrorAs(t, err, &ve) {
				assert.Equal(t, tt.expectedField, ve.Field, "Field should be 'email'")
				assert.Contains(t, ve.Error(), tt.expectedMsg)
			}

			mockRepo.AssertNotCalled(t, "GetByEmail", mock.Anything, mock.Anything)
			mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

			mockHasher.AssertNotCalled(t, "Hash", mock.Anything)
		})
	}
}

func TestRegisterUser_ShortPassword(t *testing.T) {
	ctx, mockRepo, mockHasher := testutils.SetupUserService(t)
	svc := NewUserService(mockRepo, mockHasher)

	tests := []struct {
		name          string
		email         string
		password      string
		expectedError string
		expectedField string
		minLength     int
	}{
		{
			name:          "password length 3",
			email:         "henry@email.com",
			password:      "123",
			expectedError: "password length must be >= 8",
			expectedField: "password",
		},
		{
			name:          "password is empty",
			email:         "henry@email.com",
			password:      "",
			expectedError: "password length must be >= 8",
			expectedField: "password",
		},
		{
			name:          "password with spaces only",
			email:         "henry@email.com",
			password:      " ",
			expectedError: "password length must be >= 8",
			expectedField: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.RegisterUser(ctx, tt.email, tt.password)

			assert.Nil(t, user, "User should be nil when password is invalid")
			assert.Error(t, err, "Expected error for invalid password")

			assert.Contains(t, err.Error(), tt.expectedError)

			var ve *domain.ValidationError
			if assert.ErrorAs(t, err, &ve) {
				assert.Equal(t, tt.expectedField, ve.Field, "Field should be 'password'")
				assert.Contains(t, ve.Error(), tt.expectedError)

			}

			mockRepo.AssertNotCalled(t, "GetByEmail", mock.Anything, mock.Anything)
			mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

			mockHasher.AssertNotCalled(t, "Hash", mock.Anything)
		})
	}
}

func TestRegisterUser_HashError(t *testing.T) {
	ctx, mockRepo, mockHasher := testutils.SetupUserService(t)
	svc := NewUserService(mockRepo, mockHasher)

	email := "henry@example.com"
	password := "password123"

	hashError := errors.New("bcrypt: cost out of range")

	mockRepo.On("GetByEmail", ctx, email).Return(nil, nil)
	mockHasher.On("Hash", password).Return("", hashError)

	user, err := svc.RegisterUser(ctx, email, password)

	assert.Nil(t, user, "User should be nil when hashing fails")

	assert.Error(t, err, "Expected error from hash operation")

	var ve *domain.ValidationError
	assert.False(t, errors.As(err, &ve), "Hash error should not be wrapped as ValidationError ")

	assert.True(t, errors.Is(err, hashError), "Error should contain the original hash error")

	assert.Contains(t, err.Error(), "bcrypt: cost out of range", "Error message should contain the hash error details")

	assert.Contains(t, err.Error(), "hash password:", "Error should be wrapped with context")

	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

	mockRepo.AssertCalled(t, "GetByEmail", ctx, email)

	mockHasher.AssertCalled(t, "Hash", password)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestRegisterUser_CreateError(t *testing.T) {
	ctx, mockRepo, mockHasher := testutils.SetupUserService(t)
	svc := NewUserService(mockRepo, mockHasher)

	email := "henry@email.com"
	password := "password123"
	hashedPassword := "hashed_password_123"

	createError := errors.New("database connection failed")

	mockRepo.On("GetByEmail", ctx, email).Return(nil, nil)
	mockHasher.On("Hash", password).Return(hashedPassword, nil)

	mockRepo.On("Create", ctx, mock.MatchedBy(func(user *domain.User) bool {
		return user.Email == email && user.Password == hashedPassword
	})).Return(createError)

	user, err := svc.RegisterUser(ctx, email, password)

	assert.Nil(t, user, "User should be nil when create fails")

	assert.Error(t, err, "Expected error from create operation")

	var ve *domain.ValidationError
	assert.False(t, errors.As(err, &ve), "Create error should not be wrappe as ValidationError")

	assert.True(t, errors.Is(err, createError), "Error chain should contain the original create error")

	assert.Contains(t, err.Error(), "create user", "Error should provide context about user creation")
	assert.Contains(t, err.Error(), "database connection failed", "Error should contain the original error message")

	mockRepo.AssertCalled(t, "GetByEmail", ctx, email)
	mockHasher.AssertCalled(t, "Hash", password)
	mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
