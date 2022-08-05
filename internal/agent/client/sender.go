package client

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/agent/models"
	log "github.com/sirupsen/logrus"
	"net/http"
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
			err := sender.sendRequest(buildURL)

			if err != nil {
				log.WithFields(log.Fields{
					"url": buildURL,
				}).Info("Error to send data")

				return false
			}

			log.WithFields(log.Fields{
				"url": buildURL,
			}).Info("Send metric data")

			return true
		})

		metrics.Gauges.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Gauge)
			buildURL := url + fmt.Sprintf("/update/gauge/%s/%f", metric.Name, metric.Value)
			err := sender.sendRequest(buildURL)

			if err != nil {
				log.WithFields(log.Fields{
					"url": buildURL,
				}).Info("Error to send data")

				return false
			}

			log.WithFields(log.Fields{
				"url": buildURL,
			}).Info("Send metric data")

			return true
		})
	}
}

func (sender *Sender) sendRequest(url string) error {
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	request.Header.Set("Content-Type", "text/plain")
	res, err := sender.client.Do(request)

	if err != nil {
		log.WithFields(log.Fields{}).Error(fmt.Printf("Request failed with status: %d \n", res.StatusCode))

		return err
	}

	defer res.Body.Close()

	return nil
}
