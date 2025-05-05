package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
	"github.com/AdamGriffiths31/pulseboard/internal/models"
)

func TestGetStatusCodeDistribution(t *testing.T) {
	tests := []struct {
		name              string
		mockReturn        map[string][]models.StatusCodeCount
		mockError         error
		expectedCode      int
		expectedBodyCheck func(string) bool
	}{
		{
			name: "returns status code distribution successfully",
			mockReturn: map[string][]models.StatusCodeCount{
				"https://example.com": {
					{URL: "https://example.com", StatusCode: 200, Count: 10},
					{URL: "https://example.com", StatusCode: 500, Count: 5},
				},
			},
			expectedCode: http.StatusOK,
			expectedBodyCheck: func(body string) bool {
				return len(body) > 0 && body != "[]"
			},
		},
		{
			name:         "returns empty array when no data",
			mockReturn:   map[string][]models.StatusCodeCount{},
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
				GetStatusCodeDistributionByURLFunc: func(start, end string) (map[string][]models.StatusCodeCount, error) {
					return tt.mockReturn, tt.mockError
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/statuscodedistribution?startDate=2025-01-01T00:00:00Z&endDate=2025-12-31T23:59:59Z", nil)
			rr := httptest.NewRecorder()

			handler := GetStatusCodeDistribution(mock)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rr.Code)
			}

			body := rr.Body.String()
			if !tt.expectedBodyCheck(body) {
				t.Errorf("unexpected response body: %s", body)
			}

			if tt.expectedCode == http.StatusOK && len(tt.mockReturn) > 0 {
				var decoded map[string][]models.StatusCodeCount
				if err := json.Unmarshal([]byte(body), &decoded); err != nil {
					t.Errorf("error decoding JSON: %v", err)
				}
			}
		})
	}
}
