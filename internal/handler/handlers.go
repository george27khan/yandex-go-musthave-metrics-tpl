package handler

import (
	"fmt"
	"net/http"
	"strconv"
)

type MetricHandler struct {
}

func NewMetricHandler() *MetricHandler {
	return &MetricHandler{}
}

func (h *MetricHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		valFloat float64
		valInt   int
		err      error
	)

	typ := r.PathValue("type")
	name := r.PathValue("name")
	value := r.PathValue("value")
	fmt.Println("test", typ, name, value)
	if typ != "gauge" && typ != "counter" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if typ == "gauge" {
		valFloat, err = strconv.ParseFloat(value, 64)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if typ == "counter" {
		valInt, err = strconv.Atoi(value)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	fmt.Println("test", valFloat, valInt)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
