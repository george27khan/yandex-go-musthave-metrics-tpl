package main

import (
	"net/http"
	h "yandex-go-musthave-metrics-tpl/internal/handler"
	ms "yandex-go-musthave-metrics-tpl/internal/repository/mem_storage"
	s "yandex-go-musthave-metrics-tpl/internal/service/server"
)

func main() {
	mux := http.NewServeMux()
	repository := ms.NewMemStorage()
	service := s.NewMetricService(repository)
	handler := h.NewMetricHandler(service)
	mux.HandleFunc("POST /update/{type}/{name}/{value}", handler.Update)
	mux.HandleFunc("POST /update/{type}/{value}", http.NotFound)
	mux.HandleFunc("POST /update/{type}", http.NotFound)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
