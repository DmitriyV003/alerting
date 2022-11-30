package databse

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type metricStorage struct {
	metrics  *sync.Map
	events   map[string][]chan func()
	onUpdate []func()
	pool     *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *metricStorage {
	metricStore := metricStorage{
		metrics: &sync.Map{},
		events: map[string][]chan func(){
			"OnUpdate": {
				make(chan func()),
			},
		},
		onUpdate: []func(){},
		pool:     pool,
	}
	go func() {
		chs := metricStore.events["OnUpdate"]
		for _, ch := range chs {
			for {
				handler := <-ch
				handler()
				log.Info("Event: OnUpdate")
			}
		}
	}()

	return &metricStore
}

func (s *metricStorage) AddOnUpdateListener(fn func()) {
	s.onUpdate = append(s.onUpdate, fn)
}

func (s *metricStorage) GetByNameAndType(ctx context.Context, name string, metricType string) (*model.Metric, error) {
	sql := `SELECT id, name, type, int_value, float_value FROM metrics WHERE name = $1 AND type = $2`
	var metric model.Metric

	row := s.pool.QueryRow(ctx, sql, name, metricType)
	err := row.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.IntValue, &metric.FloatValue)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error("Error to parse metric metric")

		return nil, applicationerrors.ErrInternalServer
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, applicationerrors.ErrNotFound
	}

	if metricType == model.GaugeType || metricType == model.CounterType {
		return &metric, nil
	}

	return nil, applicationerrors.ErrUnknownType
}

func (s *metricStorage) GetAll(ctx context.Context) *[]model.Metric {
	sql := `SELECT name, type, int_value, float_value FROM metrics`
	var metrics []model.Metric

	rows, err := s.pool.Query(ctx, sql)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("error to query from database")

		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var metric model.Metric
		err = rows.Scan(&metric.Name, &metric.Type, &metric.IntValue, &metric.FloatValue)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("error while iterating dataset")
		}

		metrics = append(metrics, metric)
	}

	return &metrics
}

func (s *metricStorage) Store(ctx context.Context, name string, metricType string, intValue *int64, floatValue *float64) error {
	sql := `INSERT INTO metrics (name, type, int_value, float_value) VALUES ($1, $2, $3, $4)`
	var args = make([]interface{}, 4)

	args[0] = name
	args[1] = metricType
	if intValue != nil {
		args[2] = *intValue
	} else {
		args[2] = nil
	}

	if floatValue != nil {
		args[3] = *floatValue
	} else {
		args[3] = nil
	}

	_, err := s.pool.Exec(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Error to write metric to Database")

		return applicationerrors.ErrInternalServer
	}

	return nil
}

func (s *metricStorage) Update(ctx context.Context, name string, metricType string, intValue *int64, floatValue *float64) error {
	sql := `UPDATE metrics SET int_value = $1, float_value = $2 WHERE name = $3 AND type = $4`
	var args = make([]interface{}, 4)

	if intValue != nil {
		args[0] = *intValue
	} else {
		args[0] = nil
	}

	if floatValue != nil {
		args[1] = *floatValue
	} else {
		args[1] = nil
	}
	args[2] = name
	args[3] = metricType
	_, err := s.pool.Exec(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Error to write metric to Database")

		return applicationerrors.ErrInternalServer
	}

	return nil
}

func (s *metricStorage) UpdateOrCreate(ctx context.Context, name string, value string, metricType string) error {
	metric, _ := s.GetByNameAndType(ctx, name, metricType)
	var intValue *int64
	var floatValue *float64
	var val interface{}
	var err error

	if metricType == model.GaugeType {
		val, err = strconv.ParseFloat(value, 64)
		f := val.(float64)
		floatValue = &f
	} else if metricType == model.CounterType {
		val, err = strconv.ParseInt(value, 10, 64)
		f := val.(int64)
		intValue = &f
	} else {
		return applicationerrors.ErrUnknownType
	}

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Invalid metric Value")

		return applicationerrors.ErrInvalidValue
	}

	if metric != nil {
		if metricType == model.CounterType {
			*intValue += *metric.IntValue
		}
		err = s.Update(ctx, name, metricType, intValue, floatValue)

		return err
	}

	err = s.Store(ctx, name, metricType, intValue, floatValue)
	if err != nil {
		return err
	}

	s.Emit("OnUpdate")

	return nil
}

func (s *metricStorage) DeleteByNameAndType(ctx context.Context, name string, mType string) error {
	sql := `DELETE FROM metrics WHERE name = $1 AND type = $2`
	var args = make([]interface{}, 2)

	args[0] = name
	args[1] = mType
	_, err := s.pool.Exec(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Error to delete metric from Database")

		return applicationerrors.ErrInternalServer
	}

	return nil
}

func (s *metricStorage) RestoreCollection(ctx context.Context, metrics *[]model.Metric) {
	var err error
	for _, metric := range *metrics {
		if metric.Type == model.GaugeType {
			err = s.UpdateOrCreate(ctx, metric.Name, fmt.Sprint(*metric.FloatValue), string(metric.Type))
		} else if metric.Type == model.CounterType {
			_ = s.DeleteByNameAndType(ctx, metric.Name, string(metric.Type))
			err = s.UpdateOrCreate(ctx, metric.Name, fmt.Sprint(*metric.IntValue), string(metric.Type))
		}

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Unable to save metric from file")
		}
	}
}

func (s *metricStorage) SaveCollection(ctx context.Context, metrics *[]model.Metric) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error to begin transaction")

		return applicationerrors.ErrInternalServer
	}

	for _, metric := range *metrics {
		if metric.Type == model.GaugeType {
			err = s.UpdateOrCreate(ctx, metric.Name, fmt.Sprint(*metric.FloatValue), string(metric.Type))
		} else if metric.Type == model.CounterType {
			err = s.UpdateOrCreate(ctx, metric.Name, fmt.Sprint(*metric.IntValue), string(metric.Type))
		} else {
			err = tx.Rollback(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("Unable to Rollback transaction")
				return err
			}
			return applicationerrors.ErrInternalServer
		}

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Unable to save metric collection")
			err = tx.Rollback(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("Unable to Rollback transaction")
			}
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Unable to Commit transaction")
		return err
	}

	return nil
}

func (s *metricStorage) Emit(event string) {
	if _, ok := s.events[event]; ok {
		for _, handler := range s.events[event] {
			go func(handler chan func()) {
				for _, h := range s.onUpdate {
					handler <- h
				}
			}(handler)
		}
	}
}
