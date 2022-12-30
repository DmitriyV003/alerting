package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
	"github.com/dmitriy/alerting/internal/helpers"
	"github.com/rs/zerolog/log"
)

type App struct {
	conf Config
}

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

// go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/agent/main.go

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildTime == "" {
		buildTime = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	g, gCtx := errgroup.WithContext(ctx)
	defer stop()

	log.Printf("Build version=%s, Build date=%s\n, Build commit=%s\n", buildVersion, buildTime, buildCommit)
	app := App{
		conf: Config{},
	}
	app.config()
	publicKey, err := helpers.ImportPublicKeyFromFile(app.conf.PublicKey)
	if err != nil {
		log.Warn().Err(err).Msg("error to get public key from file")
	}

	metricService := service.NewMetricService(app.conf.Key)
	pingService := service.NewPingService()
	sender := client.New(publicKey)

	log.Info().Msgf(
		"Agent starting. Poll interval: %s; Report interval: %s",
		fmt.Sprint(app.conf.PollInterval),
		fmt.Sprint(app.conf.ReportInterval),
	)

	pingService.Ping(ctx, fmt.Sprintf("http://%s/heartbeat", app.conf.Address))
	g.Go(func() error {
		metricService.GatherMetricsByInterval(gCtx, app.conf.PollInterval)
		return nil
	})
	g.Go(func() error {
		metricService.GatherAdditionalMetricsByInterval(gCtx, app.conf.PollInterval)
		return nil
	})
	g.Go(func() error {
		sender.SendWithInterval(gCtx, fmt.Sprintf("http://%s/update", app.conf.Address), &metricService.Health, app.conf.ReportInterval)
		return nil
	})

	srv2 := &http.Server{
		Addr:    ":8087",
		Handler: nil,
	}
	g.Go(func() error {
		return srv2.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return srv2.Shutdown(gCtx)
	})

	if err := g.Wait(); err != nil {
		log.Warn().Msg("agent shutdown")
	}

	log.Info().Msg("shutdown")
}
