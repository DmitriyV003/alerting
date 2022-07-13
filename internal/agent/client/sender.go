package client

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/agent/models"
	"io"
	"net/http"
	"os"
	"time"
)

type Sender struct {
	client http.Client
}

func New() Sender {
	sender := Sender{
		client: http.Client{},
	}
	sender.client.Timeout = 1 * time.Second

	return sender
}

func (sender *Sender) SendWithInterval(url string, metrics *models.Health, seconds int) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)

	for range ticker.C {
		metrics.Counters.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Counter)
			buildURL := url + fmt.Sprintf("/update/counter/%s/%d", metric.Name, metric.Value)
			sender.sendRequest(buildURL)
			fmt.Println(buildURL)

			return true
		})

		metrics.Gauges.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Gauge)
			buildURL := url + fmt.Sprintf("/update/gauge/%s/%f", metric.Name, metric.Value)
			sender.sendRequest(buildURL)
			fmt.Println(buildURL)

			return true
		})
	}
}

func (sender *Sender) sendRequest(url string) {
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	request.Header.Set("Content-Type", "text/plain")
	res, err := sender.client.Do(request)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if err != nil {
		fmt.Printf("Request failed with status: %d \n", res.StatusCode)
		os.Exit(1)
	}
}
