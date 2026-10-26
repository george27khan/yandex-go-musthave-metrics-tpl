package server

import (
	"context"
	"yandex-go-musthave-metrics-tpl/internal/model"
	storage "yandex-go-musthave-metrics-tpl/internal/repository/memstorage"
)

//go:generate mockgen -package server -source=metrics.go -destination=mock_metrics_test.go
var _ MetricRepository = (*storage.MemStorage)(nil)

type MetricRepository interface {
	Add(ctx context.Context, m model.Metrics) error
	Get(ctx context.Context, name string) (model.Metrics, error)
	GetAll(ctx context.Context) ([]model.Metrics, error)
}

type MetricService struct {
	repository MetricRepository
}

func NewMetricService(mr MetricRepository) *MetricService {
	return &MetricService{
		repository: mr,
	}
}

func (s *MetricService) Add(ctx context.Context, m model.Metrics) error {
	if m.MType == model.Counter {
		oldMetric, err := s.repository.Get(ctx, m.ID)
		if err == nil {
			*m.Delta += *oldMetric.Delta
		}
	}
	if err := s.repository.Add(ctx, m); err != nil {
		return err
	}
	return nil
}

func (s *MetricService) GetAll(ctx context.Context) ([]model.Metrics, error) {
	metrics, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (s *MetricService) GetMetric(ctx context.Context, metric model.Metrics) (model.Metrics, error) {
	metric, err := s.repository.Get(ctx, metric.ID)
	if err != nil {
		return model.Metrics{}, err
	}
	return metric, nil
}
