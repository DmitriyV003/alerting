package client

import (
	"bytes"
	"encoding/json"
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

func (sender *Sender) SendWithInterval(url string, metrics *models.Health, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		metrics.Metrics.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Metric)
			var metricValue interface{}

			if metric.Type == models.GaugeType {
				metricValue = *metric.FloatValue
			} else if metric.Type == models.CounterType {
				metricValue = *metric.IntValue
			}

			data, err := json.Marshal(metric)

			if err != nil {
				log.Error("Unknown error during json.Marshal")

				return false
			}

			err = sender.sendRequest(url, data)

			if err != nil {
				log.WithFields(log.Fields{
					"Name":  metric.Name,
					"Type":  metric.Type,
					"Value": metricValue,
				}).Error("Error to send data")

				return false
			}

			log.WithFields(log.Fields{
				"Name":  metric.Name,
				"Type":  metric.Type,
				"Value": metricValue,
			}).Info("Send metric data")

			return true
		})
	}
}

func (sender *Sender) sendRequest(url string, data []byte) error {
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")
	res, err := sender.client.Do(request)

	if err != nil {
		log.Error("Request fail", err)

		return err
	}

	log.WithFields(log.Fields{
		"StatusCode": res.StatusCode,
	}).Info("Request ended")

	defer res.Body.Close()

	return nil
}
