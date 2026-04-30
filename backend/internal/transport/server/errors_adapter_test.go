package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngryM0e/AceClub/Backend/internal/domain"
)

func TestErrorAdapter_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "validation error returns 422",
			err:        domain.NewValidationError("email", errors.New("invalid format")),
			wantStatus: http.StatusUnprocessableEntity,
			wantMsg:    "invalid email: invalid format",
		},
		{
			name:       "conflict error",
			err:        domain.NewConflictError("user", "admin"),
			wantStatus: http.StatusConflict,
			wantMsg:    "user with admin already exists",
		},
		{
			name:       "not found error",
			err:        domain.NewNotFoundError("book", "author Stephen King"),
			wantStatus: http.StatusNotFound,
			wantMsg:    "book with author Stephen King not found",
		},
		{
			name:       "unknown error",
			err:        errors.New("database connection failed"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) error {
				return tt.err
			}

			adapter := ErrorAdapter(handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			adapter.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tt.wantStatus)
			}

			var response map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if response["error"] != tt.wantMsg {
				t.Errorf("error message: got '%s', want '%s'", response["error"], tt.wantMsg)
			}
		})
	}
}
