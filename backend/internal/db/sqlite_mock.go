package db

import (
	"github.com/AdamGriffiths31/pulseboard/internal/models"
)

type MockDBClient struct {
	StoreEndpointFunc                  func(models.MonitoredEndpoint) error
	StoreMetricFunc                    func(models.Metric) error
	GetAllEndpointsFunc                func() ([]models.MonitoredEndpoint, error)
	GetAllMetricsFunc                  func(startDate, endDate string) ([]models.Metric, error)
	GetStatusCodeDistributionByURLFunc func(startDate, endDate string) (map[string][]models.StatusCodeCount, error)
	DeleteDatabaseFunc                 func() error
	CreateDatabaseFunc                 func() error
}

func (m *MockDBClient) StoreEndpoint(ep models.MonitoredEndpoint) error {
	return m.StoreEndpointFunc(ep)
}

func (m *MockDBClient) StoreMetric(metric models.Metric) error {
	return m.StoreMetricFunc(metric)
}

func (m *MockDBClient) GetAllEndpoints() ([]models.MonitoredEndpoint, error) {
	return m.GetAllEndpointsFunc()
}

func (m *MockDBClient) GetAllMetrics(startDate, endDate string) ([]models.Metric, error) {
	return m.GetAllMetricsFunc(startDate, endDate)
}

func (m *MockDBClient) GetStatusCodeDistributionByURL(startDate, endDate string) (map[string][]models.StatusCodeCount, error) {
	return m.GetStatusCodeDistributionByURLFunc(startDate, endDate)
}

func (m *MockDBClient) DeleteDatabase() error {
	return m.DeleteDatabaseFunc()
}

func (m *MockDBClient) CreateDatabase() error {
	return m.CreateDatabaseFunc()
}
