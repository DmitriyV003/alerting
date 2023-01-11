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
	GrpcAddress   string           `json:"grpc_address"`
	StoreInterval helpers.Duration `json:"store_interval"`
	StoreFile     string           `json:"store_file"`
	Restore       string           `json:"restore"`
	Key           string           `json:"key"`
	DatabaseDsn   string           `json:"database_dsn"`
	PrivateKey    string           `json:"crypto_key"`
}

type Config struct {
	Address       string        `env:"ADDRESS"`
	GrpcAddress   string        `env:"GRPC_ADDRESS"`
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
	var envConfig Config
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
	conf.GrpcAddress = jsonConfig.GrpcAddress

	storeIntervalDuration, _ := time.ParseDuration("0s")

	address := flag.String("a", "", "Server address")
	grpcAddress := flag.String("grpc-address", "", "gRPC Server address")
	storeInterval := flag.Duration("i", storeIntervalDuration, "Store data on disk interval")
	storeFile := flag.String("f", "", "File storage for data")
	restore := flag.String("r", "", "Restore data from file on restart")
	key := flag.String("k", "", "Key for hashing")
	databaseDsn := flag.String("d", "", "connection string to database")
	privateKey := flag.String("crypto-key", "", "Private key")
	flag.Parse()

	if conf.Address == "" && *address != "" {
		conf.Address = *address
	}
	if conf.GrpcAddress == "" && *grpcAddress != "" {
		conf.GrpcAddress = *grpcAddress
	}
	if conf.StoreInterval.String() == "0s" && (*storeInterval).String() != "0s" {
		conf.StoreInterval = *storeInterval
	}
	if conf.StoreFile == "" && *storeFile != "" {
		conf.StoreFile = *storeFile
	}
	if conf.Restore == "" && *restore != "" {
		conf.Restore = *restore
	}
	if conf.Key == "" && *key != "" {
		conf.Key = *key
	}
	if conf.DatabaseDsn == "" && *databaseDsn != "" {
		conf.DatabaseDsn = *databaseDsn
	}
	if conf.PrivateKey == "" && *privateKey != "" {
		conf.PrivateKey = *privateKey
	}

	err = env.Parse(&envConfig)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
	}

	if envConfig.Address != "" {
		conf.Address = envConfig.Address
	}
	if envConfig.GrpcAddress != "" {
		conf.GrpcAddress = envConfig.GrpcAddress
	}
	if envConfig.StoreInterval.String() != "0s" {
		conf.StoreInterval = envConfig.StoreInterval
	}
	if envConfig.StoreFile != "" {
		conf.StoreFile = envConfig.StoreFile
	}
	if envConfig.Restore != "" {
		conf.Restore = envConfig.Restore
	}
	if envConfig.Key != "" {
		conf.Key = envConfig.Key
	}
	if envConfig.DatabaseDsn != "" {
		conf.DatabaseDsn = envConfig.DatabaseDsn
	}
	if envConfig.PrivateKey != "" {
		conf.PrivateKey = envConfig.PrivateKey
	}

	if conf.Address == "" {
		conf.Address = defaultAddress
	}
	if conf.StoreInterval.String() == "0s" {
		storeIntervalDuration, _ := time.ParseDuration(defaultStoreInterval)
		conf.StoreInterval = storeIntervalDuration
	}
	if conf.StoreFile == "" {
		conf.StoreFile = defaultStoreFile
	}
	if conf.Restore == "" {
		conf.Restore = defaultRestore
	}
	if conf.Key == "" {
		conf.Key = defaultKey
	}
	if conf.DatabaseDsn == "" {
		conf.DatabaseDsn = defaultDatabaseDsn
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
