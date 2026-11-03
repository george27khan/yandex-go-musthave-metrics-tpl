package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"yandex-go-musthave-metrics-tpl/internal/model"
	"yandex-go-musthave-metrics-tpl/internal/service/agent"
)

func Start() {
	var (
		metrics []model.Metrics
		poolCnt int64 = 0
	)
	serverHost := flag.String("a", "localhost:8080", "Set the path to the metrics endpoint")
	reportInterval := flag.Int("r", 10, "Set report interval in seconds")
	poolInterval := flag.Int("p", 2, "Set pool interval in seconds")
	flag.Parse()

	if val, ok := os.LookupEnv("ADDRESS"); ok {
		*serverHost = val
	}
	if val, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Could not parse REPORT_INTERVAL environment variable, %v", err)
			return
		}
		*reportInterval = valInt
	}
	if val, ok := os.LookupEnv("POOL_INTERVAL"); ok {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Could not parse POOL_INTERVAL environment variable, %v", err)
			return
		}
		*poolInterval = valInt
	}

	fmt.Println("Agent listening on", *serverHost)
	tPool := time.NewTicker((time.Duration)(*poolInterval) * time.Second)
	tSend := time.NewTicker((time.Duration)(*reportInterval) * time.Second)

	for {
		select {
		case <-tPool.C:
			//fmt.Println("pool")
			metrics = agent.GetMetrics()
			poolCnt++
		case <-tSend.C:
			//fmt.Println("send")
			metrics = append(metrics, model.Metrics{
				ID:    "PollCount",
				MType: model.Counter,
				Delta: &poolCnt,
				Value: nil,
				Hash:  "",
			})
			for _, value := range metrics {
				sendJSON(*serverHost, value)
			}
			poolCnt = 0
		}

	}
}

func send(c http.Client, serverHost string, m model.Metrics) {
	urlReq := &url.URL{}
	if m.MType == model.Counter {
		urlReq = &url.URL{Scheme: "http", Host: serverHost, Path: fmt.Sprintf("/update/%s/%s/%v", m.MType, m.ID, *m.Delta)}
	} else {
		urlReq = &url.URL{Scheme: "http", Host: serverHost, Path: fmt.Sprintf("/update/%s/%s/%v", m.MType, m.ID, *m.Value)}

	}
	resp, err := c.Post(urlReq.String(), "text/plain", nil)
	if err != nil {
		log.Printf("Could not send request to server, %v/n", err)
	} else {
		resp.Body.Close()
	}
}

func sendJSON(serverHost string, m model.Metrics) {
	body := &bytes.Buffer{}
	urlReq := &url.URL{Scheme: "http", Host: serverHost, Path: fmt.Sprintf("/update/")}
	gw, err := gzip.NewWriterLevel(body, gzip.BestCompression)
	if err != nil {
		log.Printf("Could not create gzip writer, %v/n", err)
		return
	}

	err = json.NewEncoder(gw).Encode(m)
	if err != nil {
		log.Printf("Could not encode json, %v/n", err)
		_ = gw.Close()
		return
	}
	if err = gw.Close(); err != nil {
		log.Printf("Could not close gzip writer, %v/n", err)
		return
	}

	req, err := http.NewRequest("POST", urlReq.String(), body)
	if err != nil {
		log.Printf("Could not create request, %v/n", err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Could not send request to server, %v/n", err)
		return
	}
	defer resp.Body.Close()
}
