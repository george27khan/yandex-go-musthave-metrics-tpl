package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	config "yandex-go-musthave-metrics-tpl/internal/config/server"
	h "yandex-go-musthave-metrics-tpl/internal/handler"
	"yandex-go-musthave-metrics-tpl/internal/handler/middleware"
	"yandex-go-musthave-metrics-tpl/internal/repository/filestorage"
	"yandex-go-musthave-metrics-tpl/internal/repository/memstorage"
	s "yandex-go-musthave-metrics-tpl/internal/service/server"

	"github.com/go-chi/chi/v5"
)

func Run() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err.Error())
		os.Exit(1)
	}

	repository := memstorage.NewMemStorage()
	fs, err := filestorage.NewFileStorage(repository, cfg.FileStoragePath, cfg.StoreInterval)
	if err != nil {
		logger.Error("Error creating file storage", "error", err.Error())
		os.Exit(1)
	}
	service := s.NewMetricService(repository, fs, cfg.StoreInterval == 0)
	handler := h.NewMetricHandler(service)
	middlewareLogger := middleware.NewMiddleware(logger)
	r := chi.NewRouter()
	// добавление мидлвар
	router := r.With(middlewareLogger.WithLogger, middleware.WithCompression)

	router.Post("/update/{type}/{name}/{value}", handler.Update)
	router.Get("/value/{type}/{name}", handler.Get)
	router.Post("/update/", handler.UpdateJSON)
	router.Post("/value/", handler.GetJSON)
	router.Get("/", handler.GetAll)

	// загрузка данных с файла в репозиторий
	if cfg.Restore {
		err = fs.Load(context.Background())
		if err != nil {
			logger.Error("Error loading metrics from file", "error", err.Error())
		}
	}

	// запуск синхронизации файла с репозиторием
	go func() {
		err := fs.Run(context.Background())
		if err != nil {
			logger.Error("Error file storage", "error", err.Error())
		}
	}()

	logger.Info("Listening on", slog.String("address", cfg.Host))
	if err := http.ListenAndServe(cfg.Host, router); err != nil {
		panic(err)
	}
}
