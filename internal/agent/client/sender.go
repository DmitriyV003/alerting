package client

import (
	"bytes"
	"encoding/json"
	"github.com/dmitriy/alerting/internal/agent/models"
	log "github.com/sirupsen/logrus"
	"io"
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
				log.Error("Unknown error during json.Marshal: ", err)

				return false
			}

			res, err := sender.sendRequest(url, data)

			logFields := log.Fields{
				"Name":  metric.Name,
				"Type":  metric.Type,
				"Value": metricValue,
				"Hash":  metric.Hash,
			}

			if err != nil {
				log.WithFields(logFields).Error("Error to send data: ", err)

				return false
			}

			log.WithFields(logFields).WithFields(log.Fields{
				"response body": res,
			}).Info("Send metric data")

			return true
		})
	}
}

func (sender *Sender) sendRequest(url string, data []byte) (*senderResponse, error) {
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")
	res, err := sender.client.Do(request)

	if err != nil {
		log.WithFields(log.Fields{
			"StatusCode": res.StatusCode,
			"Error":      err,
		}).Error("Request fail")

		return nil, err
	}

	log.WithFields(log.Fields{
		"StatusCode": res.StatusCode,
	}).Info("Request ended")

	response := senderResponse{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Error to read body: ", err)
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Error to Unmarshal response: ", err)
		return nil, err
	}
	log.WithFields(log.Fields{
		"Response": response,
	}).Info("Got Response")

	defer res.Body.Close()

	return &response, nil
}

type senderResponse struct {
	Hash string `json:"hash,omitempty"`
}
