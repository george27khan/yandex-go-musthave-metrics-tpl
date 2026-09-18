package server

import (
	"context"
	"fmt"
	models "yandex-go-musthave-metrics-tpl/internal/model"
	storage "yandex-go-musthave-metrics-tpl/internal/repository/mem_storage"
)

//go:generate mockgen -package server -source=metrics.go -destination=mock_metrics_test.go
var _ MetricRepository = (*storage.MemStorage)(nil)

type MetricRepository interface {
	Add(ctx context.Context, m models.Metrics) error
	Get(ctx context.Context, name string) (models.Metrics, error)
}

type MeticService struct {
	repository MetricRepository
}

func NewMetricService(mr MetricRepository) *MeticService {
	return &MeticService{
		repository: mr,
	}
}

func (s *MeticService) Add(ctx context.Context, m models.Metrics) error {
	if m.MType == models.Counter {
		oldMetric, err := s.repository.Get(ctx, m.ID)
		if err == nil {
			*m.Value += *oldMetric.Value
			fmt.Println(*m.Value)
		}
	}
	if err := s.repository.Add(ctx, m); err != nil {
		return err
	}
	return nil
}
