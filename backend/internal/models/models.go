package models

import (
	"time"

	"github.com/google/uuid"
)

type MonitoredEndpoint struct {
	ID        uuid.UUID
	URL       string
	Frequency time.Duration
	Headers   map[string]string
}

type Metric struct {
	ID         uuid.UUID `json:"id"`
	EndpointID uuid.UUID `json:"endpoint_id"`
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"status_code"`
	LatencyMS  int       `json:"latency_ms"`
	URL        string    `json:"url"`
}

type StatusCodeCount struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Count      int    `json:"count"`
}
