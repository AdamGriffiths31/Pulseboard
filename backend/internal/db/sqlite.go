package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/AdamGriffiths31/pulseboard/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteClient struct {
	DB *sql.DB
}

// Initialize the database and create necessary tables
func NewSQLiteClient(path string) (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Create tables for endpoints and metrics
	schema := `
	CREATE TABLE IF NOT EXISTS monitored_endpoints (
		id TEXT PRIMARY KEY,
		url TEXT,
		frequency INTEGER,
		headers TEXT
	);

	CREATE TABLE IF NOT EXISTS api_metrics (
		id TEXT PRIMARY KEY,
		endpoint_id TEXT,
		timestamp DATETIME,
		status_code INTEGER,
		latency_ms INTEGER,
		FOREIGN KEY(endpoint_id) REFERENCES monitored_endpoints(id)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &SQLiteClient{DB: db}, nil
}

// Store an endpoint in the database
func (c *SQLiteClient) StoreEndpoint(ep models.MonitoredEndpoint) error {
	headersJSON, err := json.Marshal(ep.Headers)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec(`
		INSERT OR REPLACE INTO monitored_endpoints (id, url, frequency, headers)
		VALUES (?, ?, ?, ?)`,
		ep.ID.String(), ep.URL, int(ep.Frequency.Seconds()), string(headersJSON),
	)
	return err
}

// Store a metric in the database
func (c *SQLiteClient) StoreMetric(m models.Metric) error {
	_, err := c.DB.Exec(`
		INSERT INTO api_metrics (id, endpoint_id, timestamp, status_code, latency_ms)
		VALUES (?, ?, ?, ?, ?)`,
		m.ID.String(), m.EndpointID.String(), m.Timestamp.Format(time.RFC3339),
		m.StatusCode, m.LatencyMS,
	)
	return err
}

// Fetch all endpoints from the DB
func (c *SQLiteClient) GetAllEndpoints() ([]models.MonitoredEndpoint, error) {
	rows, err := c.DB.Query("SELECT id, url, frequency FROM monitored_endpoints")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []models.MonitoredEndpoint
	for rows.Next() {
		var ep models.MonitoredEndpoint
		var freq int
		err := rows.Scan(&ep.ID, &ep.URL, &freq)
		if err != nil {
			return nil, err
		}
		ep.Frequency = time.Duration(freq) * time.Second
		endpoints = append(endpoints, ep)
	}

	return endpoints, nil
}

// Fetch all metrics from the DB within a date range
func (c *SQLiteClient) GetAllMetrics(startDate, endDate string) ([]models.Metric, error) {
	rows, err := c.DB.Query(`
	SELECT m.id, m.endpoint_id, m.timestamp, m.status_code, m.latency_ms, e.url
	FROM api_metrics m
	JOIN monitored_endpoints e ON m.endpoint_id = e.id
	WHERE m.timestamp BETWEEN ? AND ?
	ORDER BY m.timestamp ASC
	LIMIT 100`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var m models.Metric
		var timestamp string
		err := rows.Scan(&m.ID, &m.EndpointID, &timestamp, &m.StatusCode, &m.LatencyMS, &m.URL)
		if err != nil {
			return nil, err
		}
		m.Timestamp, _ = time.Parse(time.RFC3339, timestamp) // Convert string to time
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (c *SQLiteClient) GetStatusCodeDistributionByURL(startDate, endDate string) (map[string][]models.StatusCodeCount, error) {
	rows, err := c.DB.Query(`
		SELECT e.url, m.status_code, COUNT(*) as count
		FROM api_metrics m
		JOIN monitored_endpoints e ON m.endpoint_id = e.id
		WHERE m.timestamp BETWEEN ? AND ?
		GROUP BY e.url, m.status_code
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]models.StatusCodeCount)

	for rows.Next() {
		var url string
		var statusCode, count int

		if err := rows.Scan(&url, &statusCode, &count); err != nil {
			return nil, err
		}

		result[url] = append(result[url], models.StatusCodeCount{
			URL:        url,
			StatusCode: statusCode,
			Count:      count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *SQLiteClient) DeleteDatabase() error {
	_, err := c.DB.Exec("DROP TABLE IF EXISTS monitored_endpoints")
	if err != nil {
		return err
	}
	_, err = c.DB.Exec("DROP TABLE IF EXISTS api_metrics")
	if err != nil {
		return err
	}
	return nil
}

func (c *SQLiteClient) CreateDatabase() error {
	_, err := c.DB.Exec(`
	CREATE TABLE IF NOT EXISTS monitored_endpoints (
		id TEXT PRIMARY KEY,
		url TEXT,
		frequency INTEGER,
		headers TEXT
	);

	CREATE TABLE IF NOT EXISTS api_metrics (
		id TEXT PRIMARY KEY,
		endpoint_id TEXT,
		timestamp DATETIME,
		status_code INTEGER,
		latency_ms INTEGER,
		FOREIGN KEY(endpoint_id) REFERENCES monitored_endpoints(id)
	);`)
	return err
}
