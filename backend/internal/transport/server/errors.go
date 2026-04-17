package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// Custom type with error return over http.Handler
type APIHandler func(w http.ResponseWriter, r *http.Request) error

// Adapter for convert APIHandler to http.HandleFunc
func ErrorAdapter(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// Log error details
			log.Printf("Error: %v | Path: %s | Method: %s",
				err, r.URL.Path, r.Method)

			// Send JSON with error to client
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
	}
}
