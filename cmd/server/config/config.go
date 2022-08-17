package config

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
}

const DefaultAddress = "localhost:8080"
const DefaultStoreInterval = "300s"
const DefaultStoreFile = "/tmp/devops-metrics-db.json"
const DefaultRestore = "true"

func ParseEnv(conf *Config) {
	err := env.Parse(conf)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
	}

	storeIntervalDuration, _ := time.ParseDuration(DefaultStoreInterval)

	address := flag.String("a", DefaultAddress, "Server address")
	storeInterval := flag.Duration("i", storeIntervalDuration, "Store data on disk interval")
	storeFile := flag.String("f", DefaultStoreFile, "File storage for data")
	restore := flag.String("r", DefaultRestore, "Restore data from file on restart")
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
}