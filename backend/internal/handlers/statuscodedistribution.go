package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
)

// Handler function to get the latest metrics
func GetStatusCodeDistribution(dbClient *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request for status code distribution from %s", r.RemoteAddr)
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Make this configurable
		w.Header().Set("Content-Type", "application/json")

		startDateStr := r.URL.Query().Get("startDate")
		endDateStr := r.URL.Query().Get("endDate")

		log.Printf("Start date: %s, End date: %s", startDateStr, endDateStr)

		// Fetch the latest metrics from the database
		metrics, err := dbClient.GetStatusCodeDistributionByURL(startDateStr, endDateStr)
		if err != nil {
			log.Printf("Database error while fetching status code distribution: %v", err)
			http.Error(w, "Internal server error while fetching metrics", http.StatusInternalServerError)
			return
		}

		// If no metrics, return an empty array with 200 OK
		if len(metrics) == 0 {
			log.Println("No status code distribution metrics found in the database")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("[]")); err != nil {
				log.Printf("Error writing empty response: %v", err)
			}
			return
		}

		// Convert metrics to JSON and send response
		if err := json.NewEncoder(w).Encode(metrics); err != nil {
			log.Printf("Error encoding status code distribution metrics to JSON: %v", err)
			http.Error(w, "Internal server error while encoding metrics", http.StatusInternalServerError)
		}
	}
}
