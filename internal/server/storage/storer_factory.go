package storage

import (
	"errors"

	"github.com/dmitriy/alerting/internal/server/storage/databse"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"github.com/jackc/pgx/v4/pgxpool"
)

type storageType string

type StorerFactory struct {
	Pool        *pgxpool.Pool
	FileURL     string
	DatabaseDsn string
}

func (sf *StorerFactory) Storage(st *storageType) (MetricStorage, error) {
	storerType := sf.chooseType(st)
	if storerType == "database" {
		return databse.New(sf.Pool), nil
	} else if storerType == "memory" {
		return memory.New(), nil
	}

	return nil, errors.New("unknown storage type")
}

func (sf *StorerFactory) chooseType(storageType *storageType) storageType {
	if sf.DatabaseDsn != "" || (storageType != nil && *storageType == "database") {
		return "database"
	} else if sf.FileURL != "" || (storageType != nil && *storageType == "memory") {
		return "memory"
	}

	return ""
}
