package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

const DefaultAddress = "localhost:8080"
const DefaultReportInterval = "10s"
const DefaultPollInterval = "2s"

func (conf *Config) parseEnv() {
	err := env.Parse(conf)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
	}

	reportInterval, _ := time.ParseDuration(DefaultReportInterval)
	pollInterval, _ := time.ParseDuration(DefaultPollInterval)

	address := flag.String("a", DefaultAddress, "Server address")
	reportIntervalFlag := flag.Duration("r", reportInterval, "Report Interval")
	pollIntervalFlag := flag.Duration("p", pollInterval, "Poll Interval")
	flag.PrintDefaults()
	flag.Parse()

	if conf.Address == "" {
		conf.Address = *address
	}
	if conf.ReportInterval.String() == "0s" {
		conf.ReportInterval = *reportIntervalFlag
	}
	if conf.PollInterval.String() == "0s" {
		conf.PollInterval = *pollIntervalFlag
	}
}
