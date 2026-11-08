package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
	"yandex-go-musthave-metrics-tpl/internal/model"
)

type MemStorage interface {
	GetAll(ctx context.Context) ([]model.Metrics, error)
	Add(ctx context.Context, m model.Metrics) error
}

type FileStorage struct {
	rep          MemStorage
	filePath     string
	syncInterval time.Duration
	logger       *slog.Logger
}

func NewFileStorage(rep MemStorage, filePath string, syncInterval time.Duration) (*FileStorage, error) {
	return &FileStorage{
		rep:          rep,
		filePath:     filePath,
		syncInterval: syncInterval,
	}, nil
}

// Run запуск сохранения метрик в файл по тикеру
func (fs *FileStorage) Run(ctx context.Context) error {
	ticker := time.NewTicker(fs.syncInterval)
	for {
		select {
		case <-ticker.C:
			err := fs.Save(ctx)
			if err != nil {
				fs.logger.Error("Failed to save metrics: ", slog.String("error", err.Error()))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// Save сохранение метрик с репозитория в файл
func (fs *FileStorage) Save(ctx context.Context) error {
	metrics, err := fs.rep.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshalling metrics: %w", err)
	}
	file, err := os.OpenFile(fs.filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("could not open metrics file: %w", err)
	}
	defer file.Close()
	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("writing metrics: %w", err)
	}
	return nil
}

// Load загрузка метрик с файла в репозиторий
func (fs *FileStorage) Load(ctx context.Context) error {
	file, err := os.Open(fs.filePath)
	if err != nil {
		return fmt.Errorf("could not open metrics file: %w", err)
	}
	metrics := []model.Metrics{}
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil {
		return fmt.Errorf("parsing metrics from file: %w", err)
	}
	for _, m := range metrics {
		err = fs.rep.Add(ctx, m)
		if err != nil {
			fs.logger.Error("Failed to add metric: ", slog.String("id", m.ID))
		}
	}
	return nil
}
