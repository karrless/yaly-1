package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCalcService struct{}

func (m *mockCalcService) Calc(ex string) (float64, error) {
	return 0, nil
}

func TestServerRoutes(t *testing.T) {
	mockService := &mockCalcService{}

	server := NewServer("8080", mockService)

	t.Run("Test POST /api/v1/calculate", func(t *testing.T) {

		reqBody := []byte(`{"expression":"2+2*(2+3)"}`)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(reqBody))
		rec := httptest.NewRecorder()

		server.httpServer.Handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
