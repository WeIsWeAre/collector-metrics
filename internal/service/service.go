package service

import (
	"storage"
)

type Service interface {
	GetMetric(string) (string, error)
}

type MetricService struct {
	repo storage.Repository
}

func (s *MetricService) GetMetric(longURL string) (string, error) {

	r := &storage.MetricRepository{}
	service := newGetMetricService(r)
	service.repo.Get("any")

	/*mockRepo := &MockRepository{}
	service1 := newGetMetricService(mockRepo)
	service1.repo.Get("any")*/

	return "any", nil
}

func newGetMetricService(repo storage.Repository) *MetricService {
	return &MetricService{repo: repo}
}
