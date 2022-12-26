package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/dmitriy/alerting/internal/helpers"
	log "github.com/sirupsen/logrus"
)

type JSONConfig struct {
	Address       string           `json:"address"`
	StoreInterval helpers.Duration `json:"store_interval"`
	StoreFile     string           `json:"store_file"`
	Restore       string           `json:"restore"`
	Key           string           `json:"key"`
	DatabaseDsn   string           `json:"database_dsn"`
}

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       string        `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
}

type configFile struct {
	Path string `env:"CONFIG"`
}

const defaultAddress = "localhost:8080"
const defaultStoreInterval = "300s"
const defaultStoreFile = "/tmp/devops-metrics-db.json"
const defaultRestore = "true"
const defaultKey = ""
const defaultDatabaseDsn = ""

func (conf *Config) parseEnv() {
	var jsonConfig JSONConfig
	var confFile configFile
	err := env.Parse(&confFile)
	if err != nil {
		log.Warn("Unable to parse path to config from ENV: ", err)
	}
	confPath := flag.String("config", "/home/dmitriy/GolandProjects/alerting/cmd/agent/config.json", "Config file")
	if confPath != nil && *confPath != "" {
		confFile.Path = *confPath
	}

	if confFile.Path != "" {
		err := initConfigFromJSONFile(confFile.Path, &jsonConfig)
		if err != nil {
			log.Warnf("Unable to parse configFile: %s", confFile.Path)
		}
	}

	err = env.Parse(conf)
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

	log.WithFields(log.Fields{
		"Address":       conf.Address,
		"StoreInterval": conf.StoreInterval,
		"StoreFile":     conf.StoreFile,
		"Restore":       conf.Restore,
		"Key":           conf.Key,
		"DatabaseDsn":   conf.DatabaseDsn,
	}).Info("Environment variables")
}

func initConfigFromJSONFile(file string, config *JSONConfig) error {
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonFile, config)
	log.Print("FROM FILE ", config)
	if err != nil {
		return err
	}
	return nil
}
