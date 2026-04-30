package postgres

import (
	"errors"
	"testing"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/AngryM0e/AceClub/Backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUserRepo_Create(t *testing.T) {
	ctx, db, cleanupData, fullCleanup := testutils.SetupTest(t, "../../../migrations")
	defer fullCleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name        string
		user        *domain.User
		wantErr     bool
		wantErrType error
	}{
		{
			name: "successfully creates user",
			user: &domain.User{
				Email: "henry@email.com",
			},
			wantErr: false,
		},
		{
			name: "fails with duplicate email",
			user: &domain.User{
				Email: "henry@email.com",
			},
			wantErr:     true,
			wantErrType: &domain.ConflictError{},
		},
		{
			name: "fails with empty email",
			user: &domain.User{
				Email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "fails with duplicate email" {
				cleanupData()
			}

			err := repo.Create(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "fails with duplicate email" {
					var ce *domain.ConflictError
					assert.True(t, errors.As(err, &ce))
				}
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.user.ID)
			}
		})
	}
}

func TestUserRepo_GetByEmail(t *testing.T) {
	ctx, db, cleanupData, fullCleanup := testutils.SetupTest(t, "../../../migrations")
	defer fullCleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name         string
		setupEmail   string
		searchEmail  string
		wantErr      bool
		wantNotFound bool
	}{
		{
			name:        "successfully get user by email",
			setupEmail:  "henry@email.com",
			searchEmail: "henry@email.com",
			wantErr:     false,
		},
		{
			name:         "get wrong email - not found",
			setupEmail:   "henry@email.com",
			searchEmail:  "william@email.com",
			wantErr:      true,
			wantNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupData()

			user := &domain.User{
				Email: tt.setupEmail,
			}
			err := repo.Create(ctx, user)
			require.NoError(t, err)

			found, err := repo.GetByEmail(ctx, tt.searchEmail)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.NotNil(t, found)
				assert.Equal(t, tt.searchEmail, found.Email)
				return
			}

			assert.Error(t, err)
			assert.Nil(t, found)

			if tt.wantNotFound {
				var fe *domain.NotFoundError
				require.ErrorAs(t, err, &fe)
				assert.Equal(t, "user", fe.Resource)
				assert.Equal(t, tt.searchEmail, fe.ID)
			}
		})
	}
}
