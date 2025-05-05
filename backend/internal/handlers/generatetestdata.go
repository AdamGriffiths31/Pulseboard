package handlers

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
	"github.com/AdamGriffiths31/pulseboard/internal/models"
	"github.com/google/uuid"
)

func GenerateTestData(dbClient *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Generating test data...")

		// Clear existing data
		if err := dbClient.DeleteDatabase(); err != nil {
			http.Error(w, "Failed to delete database", http.StatusInternalServerError)
			return
		}

		if err := dbClient.CreateDatabase(); err != nil {
			http.Error(w, "Failed to create database", http.StatusInternalServerError)
			return
		}

		// Insert predefined endpoints
		endpoints := []models.MonitoredEndpoint{
			{
				ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				URL:       "https://api.github.com",
				Frequency: 10 * time.Second,
				Headers:   map[string]string{"Authorization": "Bearer token"},
			},
			{
				ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
				URL:       "https://httpstat.us/200",
				Frequency: 15 * time.Second,
				Headers:   map[string]string{"Content-Type": "application/json"},
			},
			{
				ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
				URL:       "https://httpstat.us/503",
				Frequency: 20 * time.Second,
				Headers:   map[string]string{"Cache-Control": "no-cache"},
			},
		}

		for _, ep := range endpoints {
			if err := dbClient.StoreEndpoint(ep); err != nil {
				http.Error(w, "Failed to store endpoint: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Generate random metrics for each endpoint
		for _, ep := range endpoints {
			for i := 0; i < 20; i++ {
				metric := models.Metric{
					ID:         uuid.New(),
					EndpointID: ep.ID,
					Timestamp:  time.Now().Add(-24 * time.Hour).Add(time.Duration(i) * time.Hour), // Start at -24 hours and move forward
					StatusCode: []int{200, 201, 400, 401, 500, 503}[rand.Intn(6)],
					LatencyMS:  rand.Intn(1000),
				}
				if err := dbClient.StoreMetric(metric); err != nil {
					http.Error(w, "Failed to store metric: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test data generated successfully"))
	}
}
