package server

import (
	"context"
	"yandex-go-musthave-metrics-tpl/internal/model"
	storage "yandex-go-musthave-metrics-tpl/internal/repository/postgres"
)

// var _ MetricRepository = (*storage.MemStorage)(nil)
//
//go:generate mockgen -package server -source=metrics.go -destination=mock_metrics_test.go
var _ MetricRepository = (*storage.DBStorage)(nil)

type MetricRepository interface {
	Add(ctx context.Context, m model.Metrics) error
	AddList(ctx context.Context, m []model.Metrics) error
	Get(ctx context.Context, name string) (model.Metrics, error)
	GetAll(ctx context.Context) ([]model.Metrics, error)
	Ping(ctx context.Context) error
}

type MetricFile interface {
	Save(ctx context.Context) error
}

type MetricService struct {
	repository  MetricRepository
	fileStorage MetricFile
	snc         bool
}

func NewMetricService(mr MetricRepository, mf MetricFile, snc bool) *MetricService {
	return &MetricService{
		repository:  mr,
		fileStorage: mf,
		snc:         snc,
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
	if s.snc {
		return s.fileStorage.Save(ctx)
	}
	return nil
}

func (s *MetricService) AddList(ctx context.Context, metrics []model.Metrics) error {
	if err := s.repository.AddList(ctx, metrics); err != nil {
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

func (s *MetricService) PingRepository(ctx context.Context) error {
	return s.repository.Ping(ctx)
}
