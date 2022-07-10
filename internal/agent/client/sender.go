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
	ch := make(chan *http.Response)

	for range ticker.C {
		metrics.Counters.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Counter)
			buildURL := url + fmt.Sprintf("/update/counter/%s/%d", metric.Name, metric.Value)
			go sender.sendRequest(buildURL, ch)
			fmt.Println(buildURL)
			res := <-ch

			if res.StatusCode != http.StatusOK {
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(res.Body)
				fmt.Printf("Request failed with status: %d \n", res.StatusCode)
				os.Exit(1)
			}

			return true
		})

		metrics.Gauges.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Gauge)
			buildURL := url + fmt.Sprintf("/update/gauge/%s/%f", metric.Name, metric.Value)
			go sender.sendRequest(buildURL, ch)
			fmt.Println(buildURL)
			res := <-ch

			if res.StatusCode != http.StatusOK {
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(res.Body)
				fmt.Printf("Request failed with status: %d \n", res.StatusCode)
				os.Exit(1)
			}

			return true
		})
	}
}

func (sender *Sender) sendRequest(url string, ch chan *http.Response) {
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	request.Header.Set("Content-Type", "text/plain")
	res, _ := sender.client.Do(request)

	if res.StatusCode != http.StatusOK {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		fmt.Printf("Request failed with status: %d \n", res.StatusCode)
		os.Exit(1)
	}

	ch <- res
}
