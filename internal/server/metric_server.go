package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"strconv"

	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/dmitriy/alerting/internal/proto"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/dmitriy/alerting/internal/server/transformers"
	"github.com/rs/zerolog/log"
)

type MetricsServer struct {
	proto.UnimplementedMetricsServer
	store             storage.MetricStorage
	mHasher           *hasher.Hasher
	key               string
	metricTransformer *transformers.MetricTransformer
}

func NewMetricServer(store storage.MetricStorage, mHasher *hasher.Hasher, key string) *MetricsServer {
	return &MetricsServer{
		store:             store,
		mHasher:           mHasher,
		key:               key,
		metricTransformer: transformers.NewTransformer(key),
	}
}

func (ms *MetricsServer) GetMetricValueByTypeAndName(ctx context.Context, req *proto.GetMetricValueByTypeAndNameRequest) (*proto.GetMetricValueByTypeAndNameResponse, error) {
	metric, err := ms.store.GetByNameAndType(ctx, req.Name, req.Type)

	if req.Name == "" {
		log.Info().Msg("Metric Not Found ")

		return &proto.GetMetricValueByTypeAndNameResponse{}, nil
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "error to get metric: %s", err.Error())
	}

	var value float64

	if metric.Type == models.GaugeType {
		value = *metric.FloatValue
	} else if metric.Type == models.CounterType {
		value = float64(*metric.IntValue)
	}

	return &proto.GetMetricValueByTypeAndNameResponse{
		Value: value,
	}, nil
}

func (ms *MetricsServer) UpdateMetric(ctx context.Context, req *proto.UpdateMetricRequest) (*proto.UpdateMetricResponse, error) {
	var metricReq model.Metric
	var calculatedHash string
	var stringToHash string

	if metricReq.Hash != "" {
		calculatedHash = ms.mHasher.Hash(stringToHash)

		if !ms.mHasher.IsEqual(metricReq.Hash, calculatedHash) {
			return nil, status.Error(codes.DataLoss, "hashes do not equal")
		}
	}

	if req.Name == "" {
		return nil, status.Error(codes.NotFound, "name do not provide")
	}

	err := ms.store.UpdateOrCreate(ctx, req.Name, req.Value, req.Type)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error to update metric: %s", err.Error())
	}

	return &proto.UpdateMetricResponse{
		Hash: calculatedHash,
	}, nil
}

func (ms *MetricsServer) GetMetricByTypeAndName(ctx context.Context, req *proto.GetMetricByTypeAndNameRequest) (*proto.GetMetricByTypeAndNameResponse, error) {
	metric, err := ms.store.GetByNameAndType(context.Background(), req.Name, string(req.Type))

	if req.Name == "" {
		return nil, status.Error(codes.NotFound, "name do not provide")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error to GetByNameAndType: %s", err.Error())
	}

	metric = ms.metricTransformer.AddHash(metric)
	var intVal int64
	var floatVal float64
	if metric.Type == model.GaugeType {
		floatVal = *metric.FloatValue
		intVal = 0
	} else if metric.Type == model.CounterType {
		intVal = *metric.IntValue
		floatVal = 0
	}

	return &proto.GetMetricByTypeAndNameResponse{
		Name:  metric.Name,
		Type:  string(metric.Type),
		Delta: intVal,
		Value: floatVal,
		Hash:  metric.Hash,
	}, nil
}

func (ms *MetricsServer) UpdateMetricsCollection(ctx context.Context, req *proto.UpdateMetricsCollectionRequest) (*proto.UpdateMetricsCollectionResponse, error) {
	var metrics []model.Metric

	for i, metricReq := range req.Data {
		var intVal *int64
		var floatVal *float64
		if metricReq.Type == model.GaugeType {
			val, err := strconv.ParseFloat(metricReq.Value, 64)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error to parse metric value with name: %s: %s", metricReq.Name, err.Error())
			}
			floatVal = &val
		} else if metricReq.Type == model.CounterType {
			val, err := strconv.ParseInt(metricReq.Value, 10, 64)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error to parse metric value with name: %s: %s", metricReq.Name, err.Error())
			}
			intVal = &val
		}
		metric := model.Metric{
			Name:       metricReq.Name,
			Type:       model.MetricType(metricReq.Type),
			IntValue:   intVal,
			FloatValue: floatVal,
			Hash:       metricReq.Hash,
		}

		metrics[i] = metric
	}

	err := ms.store.SaveCollection(ctx, &metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error to SaveCollection: %s", err.Error())
	}

	return &proto.UpdateMetricsCollectionResponse{}, nil
}
