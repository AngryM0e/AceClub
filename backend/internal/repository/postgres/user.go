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
		return ErrEmptyEmail
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateEmail
	}

	query = "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	var lastInsertedId int64
	err = r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&lastInsertedId)
	if err != nil {
		return err
	}
	user.ID = int(lastInsertedId)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, name, email FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &ErrNotFound{Resource: "user", ID: email}
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
