package main

import (
	"net/http"
	h "yandex-go-musthave-metrics-tpl/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	handler := h.NewMetricHandler()
	mux.HandleFunc("POST /update/{type}/{name}/{value}", handler.Update)
	mux.HandleFunc("POST /update/{type}/{value}", http.NotFound)
	mux.HandleFunc("POST /update/{type}", http.NotFound)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
