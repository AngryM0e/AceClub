package domain

import "fmt"

type ValidationError struct {
	Field string
	Err   error
}

type ConflictError struct {
	Resource string
	ID       string
}

type NotFoundError struct {
	Resource string
	ID       string
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{
		Field: field,
		Err:   err,
	}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %v", e.Field, e.Err)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

func NewConflictError(resource, id string) *ConflictError {
	return &ConflictError{
		Resource: resource,
		ID:       id,
	}
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%s with %s already exists", e.Resource, e.ID)
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with %s not found", e.Resource, e.ID)
}

func (e NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}
