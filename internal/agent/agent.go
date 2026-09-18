package agent

import (
	"fmt"
	"net/http"
	"time"
	"yandex-go-musthave-metrics-tpl/internal/model"
	"yandex-go-musthave-metrics-tpl/internal/service/agent"
)

var (
	pollInterval   int = 2
	reportInterval int = 10
)

func Start() {
	var (
		metrics map[string]model.GaugeT
		poolCnt int = 0
	)
	client := http.Client{}
	tPool := time.NewTicker((time.Duration)(pollInterval) * time.Second)
	tSend := time.NewTicker((time.Duration)(reportInterval) * time.Second)

	for {
		fmt.Println("start")
		select {
		case <-tPool.C:
			fmt.Println("pool")
			metrics = agent.GetMetrics()
			poolCnt++
		case <-tSend.C:
			fmt.Println("send")
			for name, value := range metrics {
				url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", model.Gauge, name, value)
				client.Post(url, "text/plain", nil)
			}
			url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", model.Counter, "PollCount", poolCnt)
			client.Post(url, "text/plain", nil)
			poolCnt = 0
		default:
			time.Sleep((time.Duration)(pollInterval) * time.Second)
		}

	}
}
