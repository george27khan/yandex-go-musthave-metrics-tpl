package agent

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"yandex-go-musthave-metrics-tpl/internal/model"
	"yandex-go-musthave-metrics-tpl/internal/service/agent"
)

func Start() {
	var (
		metrics map[string]model.GaugeT
		poolCnt int = 0
		urlReq  *url.URL
	)
	serverHost := flag.String("a", "localhost:8080", "Set the path to the metrics endpoint")
	reportInterval := flag.Int("r", 10, "Set report interval in seconds")
	poolInterval := flag.Int("p", 2, "Set pool interval in seconds")

	flag.Parse()

	fmt.Println("Agent listening on", *serverHost)
	client := http.Client{}
	tPool := time.NewTicker((time.Duration)(*poolInterval) * time.Second)
	tSend := time.NewTicker((time.Duration)(*reportInterval) * time.Second)

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
				urlReq = &url.URL{Scheme: "http", Host: *serverHost, Path: fmt.Sprintf("/update/%s/%s/%v", model.Gauge, name, value)}
				client.Post(urlReq.String(), "text/plain", nil)
			}
			urlReq = &url.URL{Scheme: "http", Host: *serverHost, Path: fmt.Sprintf("/update/%s/%s/%v", model.Counter, "PollCount", poolCnt)}
			client.Post(urlReq.String(), "text/plain", nil)
			poolCnt = 0
		default:
			time.Sleep((time.Duration)(*poolInterval) * time.Second)
		}

	}
}
