package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
	"github.com/AngryM0e/AceClub/Backend/internal/service"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.NewValidationError("body", err)
	}

	user, err := h.service.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}
	return json.NewEncoder(w).Encode(response)
}
