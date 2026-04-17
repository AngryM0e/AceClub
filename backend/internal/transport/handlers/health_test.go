package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Проверить статус код и тело ответа

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	HealthCheck(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	expected := "{\"status\":\"up\"}\n"
	if string(body) != expected {
		t.Errorf("expected {\"status\":\"up\"\n}, got %s, ", string(body))
	}
}
