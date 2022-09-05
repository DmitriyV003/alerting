package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       string        `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
}

const defaultAddress = "localhost:8080"
const defaultStoreInterval = "300s"
const defaultStoreFile = "/tmp/devops-metrics-db.json"
const defaultRestore = "true"
const defaultKey = ""
const defaultDatabaseDsn = ""

func (conf *Config) parseEnv() {
	err := env.Parse(conf)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
	}

	storeIntervalDuration, _ := time.ParseDuration(defaultStoreInterval)

	address := flag.String("a", defaultAddress, "Server address")
	storeInterval := flag.Duration("i", storeIntervalDuration, "Store data on disk interval")
	storeFile := flag.String("f", defaultStoreFile, "File storage for data")
	restore := flag.String("r", defaultRestore, "Restore data from file on restart")
	key := flag.String("k", defaultKey, "Key for hashing")
	databaseDsn := flag.String("d", defaultDatabaseDsn, "connection string to database")
	flag.PrintDefaults()
	flag.Parse()

	if conf.Address == "" {
		conf.Address = *address
	}
	if conf.StoreInterval.String() == "0s" {
		conf.StoreInterval = *storeInterval
	}
	if conf.StoreFile == "" {
		conf.StoreFile = *storeFile
	}
	if conf.Restore == "" {
		conf.Restore = *restore
	}
	if conf.Key == "" {
		conf.Key = *key
	}
	if conf.DatabaseDsn == "" {
		conf.DatabaseDsn = *databaseDsn
	}
}
