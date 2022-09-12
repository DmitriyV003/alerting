package mocks

import (
	context "context"
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
func (m *MockMetricStorage) GetAll(ctx context.Context) *[]model.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].(*[]model.Metric)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricStorageMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricStorage)(nil).GetAll), ctx)
}

// GetByNameAndType mocks base method.
func (m *MockMetricStorage) GetByNameAndType(ctx context.Context, name, metricType string) (*model.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameAndType", ctx, name, metricType)
	ret0, _ := ret[0].(*model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameAndType indicates an expected call of GetByNameAndType.
func (mr *MockMetricStorageMockRecorder) GetByNameAndType(ctx, name, metricType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameAndType", reflect.TypeOf((*MockMetricStorage)(nil).GetByNameAndType), ctx, name, metricType)
}

// RestoreCollection mocks base method.
func (m *MockMetricStorage) RestoreCollection(ctx context.Context, metrics *[]model.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RestoreCollection", ctx, metrics)
}

// RestoreCollection indicates an expected call of RestoreCollection.
func (mr *MockMetricStorageMockRecorder) RestoreCollection(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreCollection", reflect.TypeOf((*MockMetricStorage)(nil).RestoreCollection), ctx, metrics)
}

// SaveCollection mocks base method.
func (m *MockMetricStorage) SaveCollection(ctx context.Context, metrics *[]model.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCollection", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCollection indicates an expected call of SaveCollection.
func (mr *MockMetricStorageMockRecorder) SaveCollection(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCollection", reflect.TypeOf((*MockMetricStorage)(nil).SaveCollection), ctx, metrics)
}

// UpdateOrCreate mocks base method.
func (m *MockMetricStorage) UpdateOrCreate(ctx context.Context, metric, value, metricType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrCreate", ctx, metric, value, metricType)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrCreate indicates an expected call of UpdateOrCreate.
func (mr *MockMetricStorageMockRecorder) UpdateOrCreate(ctx, metric, value, metricType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrCreate", reflect.TypeOf((*MockMetricStorage)(nil).UpdateOrCreate), ctx, metric, value, metricType)
}
