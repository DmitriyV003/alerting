package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"

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

	var listen net.Listener
	var err error
	var grpcServer *grpc.Server

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
		if app.conf.GrpcAddress == "" {
			return nil
		}
		grpcServer = grpc.NewServer()
		app.applyGrpcServices(grpcServer)
		listen, err = net.Listen("tcp", app.conf.GrpcAddress)
		return err
	})
	g.Go(func() error {
		if app.conf.GrpcAddress != "" {
			log.Info("Сервер gRPC начал работу")
			if err := grpcServer.Serve(listen); err != nil {
				log.Fatal(err, "gRPC server error")
			}
			return err
		}

		return nil
	})

	g.Go(func() error {
		return srv2.ListenAndServe()
	})
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		var err error
		<-gCtx.Done()
		err = srv2.Shutdown(gCtx)
		if app.conf.GrpcAddress != "" {
			grpcServer.Stop()
			herr := listen.Close()
			err = fmt.Errorf("listen error: %w", herr)
		}

		yerr := srv.Shutdown(gCtx)
		if yerr != nil {
			err = fmt.Errorf("srv shutdown error: %w", yerr)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Warnf("server shutdown: %v", err)
	}

	log.Info("application is shutdown")
}
