package service

import (
	"context"
	"encoding/json"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type FileSaver struct {
	file     string
	interval time.Duration
	restore  bool
	store    storage.MetricStorage
}

func NewFileSaver(file string, interval time.Duration, restore bool, store storage.MetricStorage) *FileSaver {
	return &FileSaver{
		file:     file,
		interval: interval,
		restore:  restore,
		store:    store,
	}
}

func (f *FileSaver) Restore() {
	if !f.restore || f.file == "" {
		log.Info("Loading data from file disabled")

		return
	}

	file, err := ioutil.ReadFile(f.file)
	if err != nil {
		log.Error("Unable to read data from file: ", err)

		return
	}

	var metrics []model.Metric
	err = json.Unmarshal(file, &metrics)
	if err != nil {
		log.Error("Unable to deserialize data: ", err)

		return
	}

	f.store.RestoreCollection(context.Background(), &metrics)
	log.Info("Data restored from file")
}

func (f *FileSaver) StoreAllDataWithInterval() {
	if f.file == "" {
		log.Info("Store data on disk disabled")

		return
	}
	ticker := time.NewTicker(f.interval)

	for range ticker.C {
		f.StoreAllData()
	}
}

func (f *FileSaver) StoreAllData() {
	metrics := f.store.GetAll(context.Background())

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		log.Error("Unable to serialize data: ", err)

		return
	}

	fileSplitted := strings.Split(f.file, "/")
	if len(fileSplitted) > 0 {
		fileSplitted = fileSplitted[:len(fileSplitted)-1]
	}

	directoriesPath := strings.Join(fileSplitted, "/")
	err = os.MkdirAll(directoriesPath, 0777)
	if err != nil {
		log.Error("Unable to create directories: ", err)

		return
	}

	err = ioutil.WriteFile(f.file, jsonData, 0777)
	if err != nil {
		log.Error("Unable to save data on disk: ", err)

		return
	}

	log.Info("Metrics saved on disk")
}
