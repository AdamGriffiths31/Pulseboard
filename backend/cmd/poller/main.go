package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
	"github.com/AdamGriffiths31/pulseboard/internal/handlers"
	"github.com/AdamGriffiths31/pulseboard/internal/models"
	"github.com/AdamGriffiths31/pulseboard/internal/poller"
	"github.com/AdamGriffiths31/pulseboard/internal/websocket"

	"github.com/google/uuid"
)

func main() {
	log.Println("Pulseboard Poller Starting...")

	sqlClient, err := db.NewSQLiteClient("metrics.db")
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}

	endpoints, err := sqlClient.GetAllEndpoints()
	if err != nil {
		log.Fatal("Failed to load endpoints:", err)
	}

	log.Printf("Loaded %d endpoints from the database\n", len(endpoints))

	// If no endpoints exist, add some default ones
	if len(endpoints) == 0 {
		endpoints = []models.MonitoredEndpoint{
			{
				ID:        uuid.New(),
				URL:       "https://api.github.com",
				Frequency: 30 * time.Second,
				Headers:   map[string]string{"User-Agent": "Pulseboard-Poller"},
			},
			{
				ID:        uuid.New(),
				URL:       "https://httpstat.us/503",
				Frequency: 60 * time.Second,
				Headers:   map[string]string{},
			},
			{
				ID:        uuid.New(),
				URL:       "https://httpstat.us/200?sleep=10000",
				Frequency: 10 * time.Second,
				Headers:   map[string]string{},
			},
		}

		// Store the endpoints in the DB
		for _, ep := range endpoints {
			if err := sqlClient.StoreEndpoint(ep); err != nil {
				log.Fatal("Failed to store endpoint:", err)
			}
		}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	poller.StartPolling(endpoints, sqlClient)

	// Set up HTTP routes and handlers
	http.HandleFunc("/getlatency", handlers.GetLatencyMetrics(sqlClient))
	http.HandleFunc("/statuscodedistribution", handlers.GetStatusCodeDistribution(sqlClient))
	http.HandleFunc("/generatetestdata", handlers.GenerateTestData(sqlClient))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, sqlClient)
	})

	port := ":8080"
	fmt.Printf("API Server running on http://localhost%s\n", port)

	go func() {
		log.Println("Starting HTTP server...")
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down poller...")
}
