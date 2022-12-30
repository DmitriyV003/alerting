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
	PrivateKey    string           `json:"crypto_key"`
}

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       string        `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
	PrivateKey    string        `env:"CRYPTO_KEY"`
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
	confPath := flag.String("config", "", "Config file")
	if confPath != nil && *confPath != "" {
		confFile.Path = *confPath
	}

	if confFile.Path != "" {
		err := initConfigFromJSONFile(confFile.Path, &jsonConfig)
		if err != nil {
			log.Warnf("Unable to parse configFile: %s", confFile.Path)
		}
	}

	conf.Address = jsonConfig.Address
	conf.StoreInterval = jsonConfig.StoreInterval.Duration
	conf.StoreFile = jsonConfig.StoreFile
	conf.Restore = jsonConfig.Restore
	conf.Key = jsonConfig.Key
	conf.DatabaseDsn = jsonConfig.DatabaseDsn
	conf.PrivateKey = jsonConfig.PrivateKey

	storeIntervalDuration, _ := time.ParseDuration(defaultStoreInterval)

	address := flag.String("a", defaultAddress, "Server address")
	storeInterval := flag.Duration("i", storeIntervalDuration, "Store data on disk interval")
	storeFile := flag.String("f", defaultStoreFile, "File storage for data")
	restore := flag.String("r", defaultRestore, "Restore data from file on restart")
	key := flag.String("k", defaultKey, "Key for hashing")
	databaseDsn := flag.String("d", defaultDatabaseDsn, "connection string to database")
	privateKey := flag.String("crypto-key", "", "Private key")
	flag.PrintDefaults()
	flag.Parse()

	if *address != "" {
		conf.Address = *address
	}
	if (*storeInterval).String() != "0s" {
		conf.StoreInterval = *storeInterval
	}
	if *storeFile != "" {
		conf.StoreFile = *storeFile
	}
	if *restore != "" {
		conf.Restore = *restore
	}
	if *key != "" {
		conf.Key = *key
	}
	if *databaseDsn != "" {
		conf.DatabaseDsn = *databaseDsn
	}
	if *privateKey != "" {
		conf.PrivateKey = *privateKey
	}

	err = env.Parse(conf)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
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
