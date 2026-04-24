package postgres

import (
	"context"
	"testing"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/AngryM0e/AceClub/Backend/internal/testutil"
	"github.com/stretchr/testify/assert"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Test_UserRepo(t *testing.T) {
	ctx := context.Background()
	db, cleanupData, fullCleanup := testutil.SetupTest(t, ctx)
	defer fullCleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
		errType error
	}{
		{
			name: "successfully creates user",
			user: &domain.User{
				Name:  "Henry",
				Email: "henry@email.com",
			},
			wantErr: false,
		},
		{
			name: "fails with duplicate email",
			user: &domain.User{
				Name:  "Henry Second",
				Email: "henry@email.com",
			},
			wantErr: true,
			errType: ErrDuplicateEmail,
		},
		{
			name: "fails with empty email",
			user: &domain.User{
				Name:  "No Email",
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
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.user.ID)
			}
		})
	}
}
