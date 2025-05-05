package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
	"github.com/AdamGriffiths31/pulseboard/internal/models"
	"github.com/google/uuid"
)

func TestGetLatencyMetrics(t *testing.T) {
	tests := []struct {
		name              string
		mockReturn        []models.Metric
		mockError         error
		expectedCode      int
		expectedBodyCheck func(string) bool
	}{
		{
			name: "returns metrics successfully",
			mockReturn: []models.Metric{
				{
					ID:         uuid.New(),
					EndpointID: uuid.New(),
					Timestamp:  time.Now(),
					StatusCode: 200,
					LatencyMS:  123,
					URL:        "https://example.com",
				},
			},
			expectedCode: http.StatusOK,
			expectedBodyCheck: func(body string) bool {
				return len(body) > 0 && body != "[]"
			},
		},
		{
			name:         "returns empty array when no metrics",
			mockReturn:   []models.Metric{},
			expectedCode: http.StatusOK,
			expectedBodyCheck: func(body string) bool {
				return body == "[]"
			},
		},
		{
			name:         "returns 500 on DB error",
			mockError:    errors.New("db failure"),
			expectedCode: http.StatusInternalServerError,
			expectedBodyCheck: func(body string) bool {
				return body != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &db.MockDBClient{
				GetAllMetricsFunc: func(start, end string) ([]models.Metric, error) {
					return tt.mockReturn, tt.mockError
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/metrics?startDate=2025-01-01T00:00:00Z&endDate=2025-12-31T23:59:59Z", nil)
			rr := httptest.NewRecorder()

			handler := GetLatencyMetrics(mock) 
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rr.Code)
			}

			body := rr.Body.String()
			if !tt.expectedBodyCheck(body) {
				t.Errorf("unexpected response body: %s", body)
			}

			if tt.expectedCode == http.StatusOK && len(tt.mockReturn) > 0 {
				var decoded []models.Metric
				if err := json.Unmarshal([]byte(body), &decoded); err != nil {
					t.Errorf("error decoding JSON: %v", err)
				}
			}
		})
	}
}
