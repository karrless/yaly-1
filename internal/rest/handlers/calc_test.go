package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"yaly-1/pkg/errs"
)

type mockCalcService struct {
	result float64
	err    error
}

func (m *mockCalcService) Calc(ex string) (float64, error) {
	return m.result, m.err
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		mockResult     float64
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Successful calculation",
			method:         http.MethodPost,
			requestBody:    calcRequest{Expression: "2 + 2"},
			mockResult:     4,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   calcResponse{Result: 4},
		},
		{
			name:           "Invalid expression",
			method:         http.MethodPost,
			requestBody:    calcRequest{Expression: "2 + +"},
			mockResult:     0,
			mockError:      errs.ErrExpressionNotValid,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   errorResponse{Error: errs.ErrExpressionNotValid.Error()},
		},
		{
			name:           "Internal error",
			method:         http.MethodPost,
			requestBody:    calcRequest{Expression: "1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035375169860499105765512820762454900903893289440758685084551339423045832369032229481658085593321233482747978262041447231687381771809192998812504040261841248583600+1"},
			mockResult:     0,
			mockError:      errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   errorResponse{Error: errs.ErrInternal.Error()},
		},
		{
			name:           "Invalid HTTP method",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   nil,
		},
		{
			name:           "Malformed JSON request",
			method:         http.MethodPost,
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockCalcService{
				result: tt.mockResult,
				err:    tt.mockError,
			}
			handler := NewCalcHandlers(mockService)

			var reqBody []byte
			if tt.requestBody != nil {
				if body, ok := tt.requestBody.(string); ok {
					reqBody = []byte(body)
				} else {
					reqBody, _ = json.Marshal(tt.requestBody)
				}
			}
			req := httptest.NewRequest(tt.method, "/calculate", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.Calculate(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedBody != nil {
				var actualBody, expectedBody interface{}

				err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
				if err != nil {
					t.Fatalf("Failed to decode actual response body: %v", err)
				}

				expectedBytes, _ := json.Marshal(tt.expectedBody)
				err = json.Unmarshal(expectedBytes, &expectedBody)
				if err != nil {
					t.Fatalf("Failed to decode expected response body: %v", err)
				}

				if !reflect.DeepEqual(actualBody, expectedBody) {
					t.Errorf("Expected body %+v, got %+v", expectedBody, actualBody)
				}
			}
		})
	}
}
