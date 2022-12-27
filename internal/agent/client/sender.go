package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/helpers"
	"github.com/rs/zerolog/log"
)

type Sender struct {
	client    http.Client
	publicKey *rsa.PublicKey
}

func New(publicKey *rsa.PublicKey) Sender {
	sender := Sender{
		client:    http.Client{},
		publicKey: publicKey,
	}
	sender.client.Timeout = 1 * time.Second

	return sender
}

func (sender *Sender) SendWithInterval(url string, metrics *models.Health, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		metrics.Metrics.Range(func(key, value interface{}) bool {
			metric, _ := value.(models.Metric)

			_, err := sender.sendRequest(url, metric)
			return err == nil
		})
	}
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
