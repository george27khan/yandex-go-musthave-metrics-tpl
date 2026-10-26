package memstorage

import (
	"context"
	"errors"
	"yandex-go-musthave-metrics-tpl/internal/model"
)

type MemStorage struct {
	data map[string]model.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]model.Metrics)}
}

func (s *MemStorage) Add(ctx context.Context, m model.Metrics) error {
	s.data[m.ID] = m
	return nil
}

func (s *MemStorage) GetAll(ctx context.Context) ([]model.Metrics, error) {
	if len(s.data) == 0 {
		return []model.Metrics{}, errors.New("no metrics found")
	}

	result := make([]model.Metrics, 0, len(s.data))
	for _, metric := range s.data {
		result = append(result, metric)
	}
	return result, nil
}

func (s *MemStorage) Get(ctx context.Context, name string) (model.Metrics, error) {
	metric, ok := s.data[name]
	if !ok {
		return model.Metrics{}, errors.New("metric not found")
	}
	return metric, nil
}
