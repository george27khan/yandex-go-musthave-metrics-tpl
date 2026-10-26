package agent

import (
	"math/rand"
	"runtime"
	"yandex-go-musthave-metrics-tpl/internal/model"
)

func GetMetrics() []model.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return []model.Metrics{
		{ID: "Alloc", MType: model.Gauge, Value: new(float64(memStats.Alloc))},
		{ID: "BuckHashSys", MType: model.Gauge, Value: new(float64(memStats.BuckHashSys))},
		{ID: "Frees", MType: model.Gauge, Value: new(float64(memStats.Frees))},
		{ID: "GCCPUFraction", MType: model.Gauge, Value: new(memStats.GCCPUFraction)},
		{ID: "GCSys", MType: model.Gauge, Value: new(float64(memStats.GCSys))},
		{ID: "HeapAlloc", MType: model.Gauge, Value: new(float64(memStats.HeapAlloc))},
		{ID: "HeapIdle", MType: model.Gauge, Value: new(float64(memStats.HeapIdle))},
		{ID: "HeapInuse", MType: model.Gauge, Value: new(float64(memStats.HeapInuse))},
		{ID: "HeapObjects", MType: model.Gauge, Value: new(float64(memStats.HeapObjects))},
		{ID: "HeapReleased", MType: model.Gauge, Value: new(float64(memStats.HeapReleased))},
		{ID: "HeapSys", MType: model.Gauge, Value: new(float64(memStats.HeapSys))},
		{ID: "LastGC", MType: model.Gauge, Value: new(float64(memStats.LastGC))},
		{ID: "Lookups", MType: model.Gauge, Value: new(float64(memStats.Lookups))},
		{ID: "MCacheInuse", MType: model.Gauge, Value: new(float64(memStats.MCacheInuse))},
		{ID: "MCacheSys", MType: model.Gauge, Value: new(float64(memStats.MCacheSys))},
		{ID: "MSpanInuse", MType: model.Gauge, Value: new(float64(memStats.MSpanInuse))},
		{ID: "MSpanSys", MType: model.Gauge, Value: new(float64(memStats.MSpanSys))},
		{ID: "Mallocs", MType: model.Gauge, Value: new(float64(memStats.Mallocs))},
		{ID: "NextGC", MType: model.Gauge, Value: new(float64(memStats.NextGC))},
		{ID: "NumForcedGC", MType: model.Gauge, Value: new(float64(memStats.NumForcedGC))},
		{ID: "NumGC", MType: model.Gauge, Value: new(float64(memStats.NumGC))},
		{ID: "OtherSys", MType: model.Gauge, Value: new(float64(memStats.OtherSys))},
		{ID: "PauseTotalNs", MType: model.Gauge, Value: new(float64(memStats.PauseTotalNs))},
		{ID: "StackInuse", MType: model.Gauge, Value: new(float64(memStats.StackInuse))},
		{ID: "StackSys", MType: model.Gauge, Value: new(float64(memStats.StackSys))},
		{ID: "Sys", MType: model.Gauge, Value: new(float64(memStats.Sys))},
		{ID: "TotalAlloc", MType: model.Gauge, Value: new(float64(memStats.TotalAlloc))},
		{ID: "RandomValue", MType: model.Gauge, Value: new(rand.Float64())},
	}
}
