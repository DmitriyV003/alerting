package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"golang.org/x/sync/errgroup"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type App struct {
	conf Config
	pool *pgxpool.Pool
}

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

// go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/server/main.go

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
	log.Printf("Build version=%s, Build date=%s\n, Build commit=%s\n", buildVersion, buildTime, buildCommit)
	app := App{
		conf: Config{},
	}
	app.config()

	if app.pool != nil {
		defer app.pool.Close()
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, os.Interrupt)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	log.Infof("server is starting at %s", app.conf.Address)
	srv := &http.Server{
		Addr:    app.conf.Address,
		Handler: app.routes(),
	}
	srv2 := &http.Server{
		Addr:    ":8082",
		Handler: nil,
	}
	g.Go(func() error {
		return srv2.ListenAndServe()
	})
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		_ = srv2.Shutdown(gCtx)
		return srv.Shutdown(gCtx)
	})

	if err := g.Wait(); err != nil {
		log.Error("server shutdown")
	}

	log.Info("application is shutdown")
}
