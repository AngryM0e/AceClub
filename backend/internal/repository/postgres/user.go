package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if user.Email == "" {
		return domain.NewValidationError("email", errors.New("cannot be empty"))
	}

	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var lastInsertedId int64
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&lastInsertedId)
	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.NewConflictError("user", user.Email)
		}
		return err
	}
	user.ID = int(lastInsertedId)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, email, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.NewNotFoundError("user", email)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
