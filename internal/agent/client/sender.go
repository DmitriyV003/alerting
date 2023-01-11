package client

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/helpers"
	"github.com/dmitriy/alerting/internal/proto"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/rs/zerolog/log"
)

type Sender struct {
	client       http.Client
	publicKey    *rsa.PublicKey
	metricClient proto.MetricsClient
}

func New(publicKey *rsa.PublicKey, metricClient proto.MetricsClient) Sender {
	sender := Sender{
		client:       http.Client{},
		publicKey:    publicKey,
		metricClient: metricClient,
	}
	sender.client.Timeout = 1 * time.Second

	return sender
}

func (sender *Sender) SendWithInterval(ctx context.Context, url string, metrics *models.Health, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics.Metrics.Range(func(key, value interface{}) bool {
				metric, _ := value.(models.Metric)

				if sender.metricClient != nil {
					go func() {
						err := sender.sendRequestViaGrpc(ctx, metric)
						if err != nil {
							log.Error().Err(err).Msg("error to send request via gRPC")
						}
					}()
				}

				go func() {
					_, err := sender.sendRequest(url, metric)
					if err != nil {
						log.Error().Err(err).Msg("error to send request via http")
					}
				}()

				return true
			})
		case <-ctx.Done():
			return
		}
	}
}

func (sender *Sender) sendRequestViaGrpc(ctx context.Context, metric models.Metric) error {
	var value string
	if metric.Type == model.GaugeType {
		value = fmt.Sprint(*metric.FloatValue)
	} else if metric.Type == models.CounterType {
		value = fmt.Sprint(*metric.IntValue)
	}
	_, err := sender.metricClient.UpdateMetric(ctx, &proto.UpdateMetricRequest{
		Type:  string(metric.Type),
		Name:  metric.Name,
		Value: value,
	})
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"Message": err.Error(),
		}).Msg("Error Send Metric via gRPC")
	}
	log.Info().Fields(map[string]interface{}{
		"Name":  metric.Name,
		"Type":  metric.Type,
		"Value": value,
	}).Msg("Send Metric via gRPC")

	return err
}

func (sender *Sender) sendRequest(url string, data interface{}) (*senderResponse, error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("Unknown error during json.Marshal")

		return nil, err
	}

	if sender.publicKey != nil {
		byteData, err = helpers.Encrypt(sender.publicKey, byteData)
		if err != nil {
			return nil, err
		}
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteData))
	request.Header.Add("Content-Type", "application/json")
	res, err := sender.client.Do(request)
	if err != nil {
		log.Warn().Fields(map[string]interface{}{
			"StatusCode": res.StatusCode,
			"Error":      err.Error(),
			"URL":        url,
			"Data":       data,
			"Method":     "POST",
		}).Msg("Request fail")

		return nil, err
	}

	response := senderResponse{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error to read body")
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error().Err(err).Msg("Error to Unmarshal response")
		return nil, err
	}
	log.Info().Fields(map[string]interface{}{
		"Response":   response,
		"StatusCode": res.StatusCode,
		"Method":     "POST",
		"URL":        url,
		"Data":       data,
	}).Msg("Send Metric")

	defer res.Body.Close()

	return &response, nil
}

type senderResponse struct {
	Hash string `json:"hash,omitempty"`
}
