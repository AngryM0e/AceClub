package service

import (
	"testing"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/AngryM0e/AceClub/Backend/internal/repository/postgres"
	"github.com/AngryM0e/AceClub/Backend/internal/service/hasher"
	"github.com/AngryM0e/AceClub/Backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser_Integration_Success(t *testing.T) {
	ctx, db, _, fullCleanup := testutils.SetupTest(t, "../../migrations")
	defer fullCleanup()
	repo := postgres.NewUserRepository(db)
	hasher, err := hasher.NewBcryptHasher(4)
	if err != nil {
		t.Fatalf("error with create hasher: %v", err)
	}

	userService := NewUserService(repo, hasher)

	newUser := &domain.User{
		Email:    "henry@email.com",
		Password: "password123",
	}

	user, err := userService.RegisterUser(ctx, newUser.Email, newUser.Password)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.True(t, user.ID > 0)
	assert.Equal(t, newUser.Email, user.Email)

	var dbEmail string
	err = db.QueryRowContext(ctx, "SELECT email FROM users WHERE id = $1", user.ID).Scan(&dbEmail)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Email, dbEmail)

	assert.NotEqual(t, newUser.Password, user.Password)
	assert.True(t, len(user.Password) > 30)

	_, err = userService.RegisterUser(ctx, newUser.Email, newUser.Password)

	var ce *domain.ConflictError
	assert.Error(t, err)
	assert.ErrorAs(t, err, &ce)
}
