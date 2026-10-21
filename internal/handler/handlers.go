package handler

import (
	"context"
	"encoding/json"
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
}

type MetricHandler struct {
	MetricService MetricService
}

func NewMetricHandler(ms MetricService) *MetricHandler {
	return &MetricHandler{
		MetricService: ms,
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
		return
	}
	if mType == model.Counter {
		value, err := h.checkCounter(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) Get(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("type")
	name := r.PathValue("name")
	if !h.checkType(mType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, err := h.MetricService.GetMetric(r.Context(), model.Metrics{ID: name, MType: mType})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
		return
	}
	resp, err := json.Marshal(metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MetricHandler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.MetricService.Add(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) GetJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, err = h.MetricService.GetMetric(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
