package mocks

import (
	model "github.com/dmitriy/alerting/internal/server/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMetricStorage is a mock of MetricStorage interface
type MockMetricStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetricStorageMockRecorder
}

// MockMetricStorageMockRecorder is the mock recorder for MockMetricStorage
type MockMetricStorageMockRecorder struct {
	mock *MockMetricStorage
}

// NewMockMetricStorage creates a new mock instance
func NewMockMetricStorage(ctrl *gomock.Controller) *MockMetricStorage {
	mock := &MockMetricStorage{ctrl: ctrl}
	mock.recorder = &MockMetricStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockMetricStorage) EXPECT() *MockMetricStorageMockRecorder {
	return _m.recorder
}

// UpdateMetric mocks base method
func (_m *MockMetricStorage) UpdateMetric(metric string, value string, metricType string) error {
	ret := _m.ctrl.Call(_m, "UpdateMetric", metric, value, metricType)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric
func (_mr *MockMetricStorageMockRecorder) UpdateMetric(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricStorage)(nil).UpdateMetric), arg0, arg1, arg2)
}

// GetAll mocks base method
func (_m *MockMetricStorage) GetAll() *[]model.Metric {
	ret := _m.ctrl.Call(_m, "GetAll")
	ret0, _ := ret[0].(*[]model.Metric)
	return ret0
}

// GetAll indicates an expected call of GetAll
func (_mr *MockMetricStorageMockRecorder) GetAll() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetAll", reflect.TypeOf((*MockMetricStorage)(nil).GetAll))
}

// GetByNameAndType mocks base method
func (_m *MockMetricStorage) GetByNameAndType(name string, metricType string) (*model.Metric, error) {
	ret := _m.ctrl.Call(_m, "GetByNameAndType", name, metricType)
	ret0, _ := ret[0].(*model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameAndType indicates an expected call of GetByNameAndType
func (_mr *MockMetricStorageMockRecorder) GetByNameAndType(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetByNameAndType", reflect.TypeOf((*MockMetricStorage)(nil).GetByNameAndType), arg0, arg1)
}

// SaveAllMetricsData mocks base method
func (_m *MockMetricStorage) SaveAllMetricsData(metrics *[]model.Metric) {
	_m.ctrl.Call(_m, "SaveAllMetricsData", metrics)
}

// SaveAllMetricsData indicates an expected call of SaveAllMetricsData
func (_mr *MockMetricStorageMockRecorder) SaveAllMetricsData(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SaveAllMetricsData", reflect.TypeOf((*MockMetricStorage)(nil).SaveAllMetricsData), arg0)
}

// AddOnUpdateListener mocks base method
func (_m *MockMetricStorage) AddOnUpdateListener(fn func()) {
	_m.ctrl.Call(_m, "AddOnUpdateListener", fn)
}

// AddOnUpdateListener indicates an expected call of AddOnUpdateListener
func (_mr *MockMetricStorageMockRecorder) AddOnUpdateListener(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddOnUpdateListener", reflect.TypeOf((*MockMetricStorage)(nil).AddOnUpdateListener), arg0)
}
