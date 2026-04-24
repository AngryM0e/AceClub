package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/AngryM0e/AceClub/Backend/internal/repository/postgres"
)

type UserHandler struct {
	userRepo domain.UserRepository
}

func NewUserHandler(userRepo domain.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return NewAPIError("Invalid request body:"+err.Error(), http.StatusBadRequest)
	}

	if req.Name == "" {
		return NewAPIError("Name is required", http.StatusBadRequest)
	}
	if req.Email == "" {
		return NewAPIError("Email is required", http.StatusBadRequest)
	}

	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		switch {
		case errors.Is(err, postgres.ErrEmptyEmail):
			return NewAPIError("Email cannot be empty", http.StatusBadRequest)
		case errors.Is(err, postgres.ErrDuplicateEmail):
			return NewAPIError("User with this email already exists", http.StatusConflict)
		default:
			return NewAPIError("Failed to create user:"+err.Error(), http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	}
	return json.NewEncoder(w).Encode(response)
}
