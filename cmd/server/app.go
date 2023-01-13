package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/dmitriy/alerting/internal/server/service"

	"github.com/dmitriy/alerting/internal/helpers"
	"github.com/dmitriy/alerting/internal/proto"
	"github.com/dmitriy/alerting/internal/server"
	middleware2 "github.com/dmitriy/alerting/internal/server/middleware"

	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/dmitriy/alerting/internal/server/handlers"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func (app *App) routes() http.Handler {
	router := chi.NewRouter()

	app.pool = app.connectToDB()

	if app.conf.DatabaseDsn != "" && app.pool != nil {
		app.migrate()
	}

	privateKey, err := helpers.ImportPrivateKeyFromFile(app.conf.PrivateKey)
	if err != nil {
		log.Error(fmt.Errorf("error to get private key from file: %w", err))
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware2.TrustedNet(app.conf.TrustedNet))
	router.Use(middleware.Compress(5))
	router.Use(middleware2.Decrypt(privateKey))
	router.Use(middleware.Heartbeat("/heartbeat"))

	storerFactory := storage.StorerFactory{
		Pool:        app.pool,
		FileURL:     app.conf.StoreFile,
		DatabaseDsn: app.conf.DatabaseDsn,
	}

	store, err := storerFactory.Storage(nil)
	if err != nil {
		log.Fatal("unknown storage type: ", err)
	}

	app.restoreData(store)

	mHasher := hasher.New(app.conf.Key)
	updateMetricHandler := handlers.NewUpdateMetricHandler(store, mHasher)
	getAllMetricsHandler := handlers.NewGetAllMetricHandler(store)
	getMetricValueByTypeAndNameHandler := handlers.NewGetMetricValueByTypeAndNameHandler(store)
	getMetricByTypeAndNameHandler := handlers.NewGetMetricByTypeAndNameHandler(store, app.conf.Key)
	pingHandler := handlers.NewPingHandler(app.pool, context.Background())
	updateMetricsCollectionHandler := handlers.NewUpdateMetricsCollectionHandler(store)

	router.Get("/", getAllMetricsHandler.Handle)
	router.Get("/value/{type}/{name}", getMetricValueByTypeAndNameHandler.Handle)
	router.Post("/update/{type}/{name}/{value}", updateMetricHandler.Handle)
	router.Post("/value", getMetricByTypeAndNameHandler.Handle)
	router.Post("/update", updateMetricHandler.Handle)
	router.Post("/updates", updateMetricsCollectionHandler.Handle)

	router.Get("/ping", pingHandler.Handle)

	return router
}

func (app *App) restoreData(store storage.MetricStorage) {
	restore, _ := strconv.ParseBool(app.conf.Restore)
	fileSaver := service.NewFileSaver(app.conf.StoreFile, app.conf.StoreInterval, restore, store)
	fileSaver.Restore()

	if app.conf.StoreInterval.String() == "0s" {
		store.AddOnUpdateListener(fileSaver.StoreAllData)
	} else {
		go fileSaver.StoreAllDataWithInterval()
	}
}

func (app *App) applyGrpcServices(grpcServer *grpc.Server) {
	mHasher := hasher.New(app.conf.Key)
	storerFactory := storage.StorerFactory{
		Pool:        app.pool,
		FileURL:     app.conf.StoreFile,
		DatabaseDsn: app.conf.DatabaseDsn,
	}
	store, err := storerFactory.Storage(nil)
	if err != nil {
		log.Fatal("unknown storage type: ", err)
	}
	proto.RegisterMetricsServer(grpcServer, server.NewMetricServer(store, mHasher, app.conf.Key))
}

func (app *App) config() {
	app.conf.parseEnv()
}

func (app *App) connectToDB() (pool *pgxpool.Pool) {
	if app.conf.DatabaseDsn == "" {
		log.Info("Database URl not provided")
		return nil
	}

	var err error
	conf, err := pgxpool.ParseConfig(app.conf.DatabaseDsn)
	if err != nil {
		log.Error("Unable to parse Database config: ", err)
		return
	}
	conf.ConnConfig.LogLevel = 5
	conf.ConnConfig.Logger = logrusadapter.NewLogger(log.New())
	pool, err = pgxpool.ConnectConfig(context.Background(), conf)

	if err != nil {
		log.Error("Unable to connect to database: ", err)
		return
	}

	return pool
}

func (app *App) migrate() {

	//parsedDbUrl, _ := url.Parse(app.conf.DatabaseDsn)
	//cmd := exec.Command("tern", "migrate", "--migrations", "./migrations")
	//cmd.Env = append(cmd.Env, fmt.Sprintf("DATABASE=%s", strings.Trim(parsedDbUrl.Path, "/")))
	//cmd.Env = append(cmd.Env, fmt.Sprintf("DATABASE_DSN=%s", app.conf.DatabaseDsn))
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	log.Error("Error during migrations: ", err)
	//	return
	//}
	//
	//log.Info("Migrating: ", string(out))

	sql := `CREATE TABLE IF NOT EXISTS metrics(
    	id serial PRIMARY KEY,
    	name VARCHAR (255) NOT NULL,
    	type VARCHAR (255) NOT NULL,
    	int_value BIGINT,
    	float_value DOUBLE PRECISION
	)`
	_, err := app.pool.Query(context.Background(), sql)
	if err != nil {
		log.Error("Error during migration: ", err)
		return
	}

	log.Info("Migrating: ", sql)
}
