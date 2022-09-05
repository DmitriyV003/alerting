package mocks

import (
	reflect "reflect"

	model "github.com/dmitriy/alerting/internal/server/model"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricStorage is a mock of MetricStorage interface.
type MockMetricStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetricStorageMockRecorder
}

// MockMetricStorageMockRecorder is the mock recorder for MockMetricStorage.
type MockMetricStorageMockRecorder struct {
	mock *MockMetricStorage
}

// NewMockMetricStorage creates a new mock instance.
func NewMockMetricStorage(ctrl *gomock.Controller) *MockMetricStorage {
	mock := &MockMetricStorage{ctrl: ctrl}
	mock.recorder = &MockMetricStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricStorage) EXPECT() *MockMetricStorageMockRecorder {
	return m.recorder
}

// AddOnUpdateListener mocks base method.
func (m *MockMetricStorage) AddOnUpdateListener(fn func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddOnUpdateListener", fn)
}

// AddOnUpdateListener indicates an expected call of AddOnUpdateListener.
func (mr *MockMetricStorageMockRecorder) AddOnUpdateListener(fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOnUpdateListener", reflect.TypeOf((*MockMetricStorage)(nil).AddOnUpdateListener), fn)
}

// Emit mocks base method.
func (m *MockMetricStorage) Emit(event string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Emit", event)
}

// Emit indicates an expected call of Emit.
func (mr *MockMetricStorageMockRecorder) Emit(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Emit", reflect.TypeOf((*MockMetricStorage)(nil).Emit), event)
}

// GetAll mocks base method.
func (m *MockMetricStorage) GetAll() *[]model.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(*[]model.Metric)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricStorage)(nil).GetAll))
}

// GetByNameAndType mocks base method.
func (m *MockMetricStorage) GetByNameAndType(name, metricType string) (*model.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameAndType", name, metricType)
	ret0, _ := ret[0].(*model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameAndType indicates an expected call of GetByNameAndType.
func (mr *MockMetricStorageMockRecorder) GetByNameAndType(name, metricType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameAndType", reflect.TypeOf((*MockMetricStorage)(nil).GetByNameAndType), name, metricType)
}

// SaveAllMetricsData mocks base method.
func (m *MockMetricStorage) SaveAllMetricsData(metrics *[]model.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveAllMetricsData", metrics)
}

// SaveAllMetricsData indicates an expected call of SaveAllMetricsData.
func (mr *MockMetricStorageMockRecorder) SaveAllMetricsData(metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAllMetricsData", reflect.TypeOf((*MockMetricStorage)(nil).SaveAllMetricsData), metrics)
}

// UpdateMetric mocks base method.
func (m *MockMetricStorage) UpdateMetric(metric, value, metricType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", metric, value, metricType)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricStorageMockRecorder) UpdateMetric(metric, value, metricType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricStorage)(nil).UpdateMetric), metric, value, metricType)
}
