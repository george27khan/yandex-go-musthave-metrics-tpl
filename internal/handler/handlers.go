package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"yandex-go-musthave-metrics-tpl/internal/model"
	"yandex-go-musthave-metrics-tpl/internal/service/server"
)

var _ MetricService = (*server.MetricService)(nil)

type MetricService interface {
	Add(ctx context.Context, m model.Metrics) error
	GetAll(ctx context.Context) ([]model.Metrics, error)
	GetMetric(ctx context.Context, metric model.Metrics) (model.Metrics, error)
	PingRepository(ctx context.Context) error
}

type MetricHandler struct {
	MetricService MetricService
	logger        *slog.Logger
}

func NewMetricHandler(ms MetricService, l *slog.Logger) *MetricHandler {
	return &MetricHandler{
		MetricService: ms,
		logger:        l,
	}
}

func (h *MetricHandler) checkType(mType string) bool {
	if mType != model.Gauge && mType != model.Counter {
		return false
	}
	return true
}

func (h *MetricHandler) checkGauge(mValue string) (float64, error) {
	val, err := strconv.ParseFloat(mValue, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (h *MetricHandler) checkCounter(mValue string) (int64, error) {
	val, err := strconv.ParseInt(mValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (h *MetricHandler) Update(w http.ResponseWriter, r *http.Request) {
	var metric model.Metrics
	mType := r.PathValue("type")
	name := r.PathValue("name")
	value := r.PathValue("value")
	if !h.checkType(mType) {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("h.checkType(mType)")

		return
	}
	if mType == model.Counter {
		value, err := h.checkCounter(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("h.checkCounter", "error", err.Error())

			return
		}
		metric = model.Metrics{
			ID:    name,
			MType: mType,
			Delta: &value,
			Value: nil,
			Hash:  "",
		}
	}
	if mType == model.Gauge {
		value, err := h.checkGauge(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("h.checkGauge", "error", err.Error())

			return
		}
		metric = model.Metrics{
			ID:    name,
			MType: mType,
			Delta: nil,
			Value: &value,
			Hash:  "",
		}
	}

	err := h.MetricService.Add(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("h.MetricService.Add", "error", err.Error())

		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) Get(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("type")
	name := r.PathValue("name")
	if !h.checkType(mType) {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("h.checkType(mType)")

		return
	}
	metric, err := h.MetricService.GetMetric(r.Context(), model.Metrics{ID: name, MType: mType})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Error("h.MetricService.GetMetric", "error", err.Error())

		return
	}
	w.Header().Set("Content-Type", "text/plain")

	if metric.MType == model.Counter {
		_, _ = w.Write([]byte(strconv.FormatInt(*metric.Delta, 10)))
	} else {
		_, _ = w.Write([]byte(strconv.FormatFloat(*metric.Value, 'f', -1, 64)))
	}

}

func (h *MetricHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.MetricService.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Error("h.MetricService.GetAll(r.Context())", "error", err.Error())

		return
	}
	resp, err := json.Marshal(metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("json.Marshal(metrics)", "error", err.Error())

		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("w.Write(resp)", "error", err.Error())

		return
	}
}

func (h *MetricHandler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("son.NewDecoder(r.Body).Decode(&metric)", "error", err.Error())

		return
	}
	err = h.MetricService.Add(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("h.MetricService.Add", "error", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) GetJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("json.NewDecoder(r.Body).Decode(&metric)", "error", err.Error())

		return
	}
	metric, err = h.MetricService.GetMetric(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Error("h.MetricService.GetMetric", "error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("json.NewEncoder(w).Encode", "error", err.Error())
		return
	}
}

func (h *MetricHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if ok := h.MetricService.PingRepository(r.Context()); ok != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("PingRepository")
		return
	}
	w.WriteHeader(http.StatusOK)
}
