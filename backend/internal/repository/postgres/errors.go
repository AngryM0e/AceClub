package postgres

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateEmail = errors.New("user with this email already exists")
	ErrEmptyEmail     = errors.New("email can't be empty")
)

type ErrNotFound struct {
	Resource string
	ID       string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.Resource, e.ID)
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*ErrNotFound)
	return ok
}
