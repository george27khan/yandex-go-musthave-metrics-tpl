package handler

import (
	"context"
	"net/http"
	"strconv"
	"yandex-go-musthave-metrics-tpl/internal/model"
	"yandex-go-musthave-metrics-tpl/internal/service/server"
)

var _ MetricService = (*server.MeticService)(nil)

type MetricService interface {
	Add(ctx context.Context, m model.Metrics) error
}

type MetricHandler struct {
	MetricService MetricService
}

func NewMetricHandler(ms MetricService) *MetricHandler {
	return &MetricHandler{
		MetricService: ms,
	}
}

func (h *MetricHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		valFloat float64
		err      error
	)

	typ := r.PathValue("type")
	name := r.PathValue("name")
	value := r.PathValue("value")
	if typ != model.Gauge && typ != model.Counter {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valFloat, err = strconv.ParseFloat(value, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.MetricService.Add(r.Context(), model.Metrics{
		ID:    name,
		MType: typ,
		Delta: nil,
		Value: &valFloat,
		Hash:  "",
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
