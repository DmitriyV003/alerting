package service

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type PingService struct {
}

func NewPingService() *PingService {
	return &PingService{}
}

func (s *PingService) Ping(address string) {
	ticker := time.NewTicker(time.Second)
	clientPing := http.Client{}
	for range ticker.C {
		err := s.send(&clientPing, address)
		if err != nil {
			log.Warn().Err(err).Msg("Ping Server")
			continue
		}
		log.Info().Msg("Connected to server")
		ticker.Stop()
		break
	}
}

func (s *PingService) send(client *http.Client, address string) error {
	request, _ := http.NewRequest(http.MethodGet, address, nil)
	res, err := client.Do(request)
	if err != nil {
		log.Warn().Err(err).Msg("Ping Server")
		return err
	}
	defer res.Body.Close()
	return nil
}
