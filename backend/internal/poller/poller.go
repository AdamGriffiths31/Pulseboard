package poller

import (
	"log"
	"net/http"
	"time"

	"github.com/AdamGriffiths31/pulseboard/internal/db"
	"github.com/AdamGriffiths31/pulseboard/internal/models"

	"github.com/google/uuid"
)

func StartPolling(endpoints []models.MonitoredEndpoint, dbClient *db.SQLiteClient) {
	for _, ep := range endpoints {
		go func(e models.MonitoredEndpoint) {
			ticker := time.NewTicker(e.Frequency)
			defer ticker.Stop()

			for {
				<-ticker.C
				metric := checkEndpoint(e)
				log.Printf("%s | %d | %dms\n", e.URL, metric.StatusCode, metric.LatencyMS)

				if err := dbClient.StoreMetric(metric); err != nil {
					log.Println("DB error:", err)
				}
			}
		}(ep)
	}
}

func checkEndpoint(ep models.MonitoredEndpoint) models.Metric {
	start := time.Now()

	req, err := http.NewRequest("GET", ep.URL, nil) // TODO: Add support for POST/PUT if needed
	if err != nil {
		// Handle error gracefully, maybe return a failed metric
		return models.Metric{}
	}

	// Set headers if any
	for k, v := range ep.Headers {
		req.Header.Set(k, v)
	}

	const requestTimeout = 5 * time.Second
	client := http.Client{Timeout: requestTimeout}

	resp, err := client.Do(req)
	duration := time.Since(start).Milliseconds()

	status := 0
	if err == nil {
		status = resp.StatusCode
		if resp.Body != nil {
			defer resp.Body.Close()
		}
	}

	return models.Metric{
		ID:         uuid.New(),
		EndpointID: ep.ID,
		Timestamp:  time.Now(),
		StatusCode: status,
		LatencyMS:  int(duration),
	}
}
