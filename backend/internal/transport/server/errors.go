package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AngryM0e/AceClub/Backend/internal/transport/handlers"
)

// Custom type with error return over http.Handler
type APIHandler func(w http.ResponseWriter, r *http.Request) error

// Adapter for convert APIHandler to http.HandleFunc
func ErrorAdapter(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.Header().Set("Content-Type", "application/json")
			// Log error details
			status := http.StatusInternalServerError
			if apiErr, ok := err.(*handlers.APIError); ok {
				status = apiErr.StatusCode
			}
			log.Printf("Error: %v | Path: %s | Method: %s",
				err, r.URL.Path, r.Method)

			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
	}
}
