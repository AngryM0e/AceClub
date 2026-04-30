package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
)

// Custom type with error return over http.Handler
type APIHandler func(w http.ResponseWriter, r *http.Request) error

// Adapter for convert APIHandler to http.HandleFunc
func ErrorAdapter(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.Header().Set("Content-Type", "application/json")
			// Log error details
			status := statusCodeFromError(err)

			slog.Error("request failed",
				slog.String("error", err.Error()),
				slog.Int("status", status),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
			)

			msg := err.Error()
			if status == http.StatusInternalServerError {
				msg = "internal server error"
			}

			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{
				"error": msg,
			})
		}
	}
}

func statusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var (
		validationErr *domain.ValidationError
		conflictErr   *domain.ConflictError
		notFoundErr   *domain.NotFoundError
	)

	switch {
	case errors.As(err, &validationErr):
		return http.StatusUnprocessableEntity // 422
	case errors.As(err, &conflictErr):
		return http.StatusConflict // 409
	case errors.As(err, &notFoundErr):
		return http.StatusNotFound // 404
	default:
		return http.StatusInternalServerError // 500
	}
}
