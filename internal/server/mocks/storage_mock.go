package mocks

import (
	reflect "reflect"

	storage "github.com/dmitriy/alerting/internal/server/storage"
	"github.com/golang/mock/gomock"
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

// GetAll mocks base method.
func (m *MockMetricStorage) GetAll() *[]storage.MetricData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(*[]storage.MetricData)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricStorage)(nil).GetAll))
}

// GetByNameAndType mocks base method.
func (m *MockMetricStorage) GetByNameAndType(arg0, arg1 string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameAndType", arg0, arg1)
	ret0, _ := ret[0]
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameAndType indicates an expected call of GetByNameAndType.
func (mr *MockMetricStorageMockRecorder) GetByNameAndType(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameAndType", reflect.TypeOf((*MockMetricStorage)(nil).GetByNameAndType), arg0, arg1)
}

// UpdateMetric mocks base method.
func (m *MockMetricStorage) UpdateMetric(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricStorageMockRecorder) UpdateMetric(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricStorage)(nil).UpdateMetric), arg0, arg1, arg2)
}
