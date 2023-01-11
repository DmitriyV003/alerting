package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/dmitriy/alerting/internal/helpers"
	"github.com/rs/zerolog/log"
)

type JSONConfig struct {
	Address        string           `json:"address"`
	GrpcAddress    string           `json:"grpc_address"`
	ReportInterval helpers.Duration `json:"report_interval"`
	PollInterval   helpers.Duration `json:"poll_interval"`
	Key            string           `json:"key"`
	PublicKey      string           `json:"crypto_key"`
}

type Config struct {
	Address        string        `env:"ADDRESS"`
	GrpcAddress    string        `env:"GRPC_ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	PublicKey      string        `env:"CRYPTO_KEY"`
}

type configFile struct {
	Path string `env:"CONFIG"`
}

const defaultAddress = "localhost:8080"
const defaultKey = ""
const defaultReportInterval = "10s"
const defaultPollInterval = "2s"

func (conf *Config) parseEnv() {
	var jsonConfig JSONConfig
	var confFile configFile
	err := env.Parse(&confFile)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to parse path to config from ENV")
	}
	confPath := flag.String("config", "", "Config file")
	if confPath != nil && *confPath != "" {
		confFile.Path = *confPath
	}

	if confFile.Path != "" {
		err := initConfigFromJSONFile(confFile.Path, &jsonConfig)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to parse configFile: %s", confFile.Path)
		}
	}

	conf.Key = jsonConfig.Key
	conf.Address = jsonConfig.Address
	conf.GrpcAddress = jsonConfig.GrpcAddress
	conf.ReportInterval = jsonConfig.ReportInterval.Duration
	conf.PollInterval = jsonConfig.PollInterval.Duration
	conf.PublicKey = jsonConfig.PublicKey

	reportInterval, err := time.ParseDuration(defaultReportInterval)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse default interval: reportInterval")
	}

	pollInterval, err := time.ParseDuration(defaultPollInterval)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse default interval: pollInterval")
	}

	address := flag.String("a", defaultAddress, "Server address")
	key := flag.String("k", defaultKey, "Key for hashing")
	reportIntervalFlag := flag.Duration("r", reportInterval, "Report Interval")
	pollIntervalFlag := flag.Duration("p", pollInterval, "Poll Interval")
	publicKey := flag.String("crypto-key", "", "Public key")
	grpcAddress := flag.String("grpc-address", "", "gRPC address")

	flag.PrintDefaults()
	flag.Parse()

	if *address != "" {
		conf.Address = *address
	}
	if *grpcAddress != "" {
		conf.GrpcAddress = *grpcAddress
	}
	if *key != "" {
		conf.Key = *key
	}
	if (*reportIntervalFlag).String() != "0s" {
		conf.ReportInterval = *reportIntervalFlag
	}
	if (*pollIntervalFlag).String() != "0s" {
		conf.PollInterval = *pollIntervalFlag
	}
	if *publicKey != "" {
		conf.PublicKey = *publicKey
	}

	err = env.Parse(conf)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to parse ENV")
	}
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
