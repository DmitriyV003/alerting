package main

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

const defaultAddress = "localhost:8080"
const defaultKey = ""
const defaultReportInterval = "10s"
const defaultPollInterval = "2s"

func (conf *Config) parseEnv() {
	err := env.Parse(conf)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse ENV")
	}

	reportInterval, _ := time.ParseDuration(defaultReportInterval)
	pollInterval, _ := time.ParseDuration(defaultPollInterval)

	address := flag.String("a", defaultAddress, "Server address")
	key := flag.String("k", defaultKey, "Key for hashing")
	reportIntervalFlag := flag.Duration("r", reportInterval, "Report Interval")
	pollIntervalFlag := flag.Duration("p", pollInterval, "Poll Interval")
	flag.PrintDefaults()
	flag.Parse()

	if conf.Address == "" {
		conf.Address = *address
	}
	if conf.Key == "" {
		conf.Key = *key
	}
	if conf.ReportInterval.String() == "0s" {
		conf.ReportInterval = *reportIntervalFlag
	}
	if conf.PollInterval.String() == "0s" {
		conf.PollInterval = *pollIntervalFlag
	}
}
