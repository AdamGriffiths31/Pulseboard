package websocket

import (
	"log"
	"net/http"
	"time"

	"math/rand"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
	"github.com/AdamGriffiths31/pulseboard/internal/models"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, dbClient *db.SQLiteClient) {
	log.Printf("Received WebSocket connection from %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	ticker := time.NewTicker(5 * time.Second) // Send updates every 5 seconds
	defer ticker.Stop()

	for range ticker.C {
		data := getTestData() // Generate fake data
		if err := conn.WriteJSON(data); err != nil {
			log.Println("Error sending data over WebSocket:", err)
			return
		}
	}
}

func getTestData() []models.Metric {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator using math/rand
	return []models.Metric{
		{
			ID:         uuid.New(),
			EndpointID: uuid.New(),
			Timestamp:  time.Now(),
			StatusCode: 200,
			LatencyMS:  rand.Intn(500), // Random latency between 0 and 499 ms
			URL:        "https://api.github.com",
		},
		{
			ID:         uuid.New(),
			EndpointID: uuid.New(),
			Timestamp:  time.Now(),
			StatusCode: 200,
			LatencyMS:  rand.Intn(500), // Random latency between 0 and 499 ms
			URL:        "https://httpstat.us/200",
		},
		{
			ID:         uuid.New(),
			EndpointID: uuid.New(),
			Timestamp:  time.Now(),
			StatusCode: 503,
			LatencyMS:  rand.Intn(500), // Random latency between 0 and 499 ms
			URL:        "https://httpstat.us/503",
		},
	}
}
