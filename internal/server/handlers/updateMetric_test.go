package handlers

import (
	"errors"
	"fmt"
	"github.com/dmitriy/alerting/internal/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler_Handle(t *testing.T) {

	tests := []struct {
		name string
		args struct {
			_type string
			name  string
			value string
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test successfull update counter",
			args: struct {
				_type string
				name  string
				value string
			}{_type: "counter", name: "testCounter", value: "100"},
			want: struct {
				statusCode int
			}{statusCode: 200},
		},
		{
			name: "test successfull update gauge",
			args: struct {
				_type string
				name  string
				value string
			}{_type: "gauge", name: "testGauge", value: "0.2222"},
			want: struct{ statusCode int }{statusCode: 200},
		},
		{
			name: "test update metric without id",
			args: struct {
				_type string
				name  string
				value string
			}{_type: "counter", name: "", value: "0.2222"},
			want: struct{ statusCode int }{statusCode: 404},
		},
		{
			name: "test update metric with incorrect type",
			args: struct {
				_type string
				name  string
				value string
			}{_type: "counterTest", name: "counterInc", value: "20"},
			want: struct{ statusCode int }{statusCode: 501},
		},
		{
			name: "test update metric with invalid value",
			args: struct {
				_type string
				name  string
				value string
			}{_type: "counter", name: "counterInc", value: "vxfg"},
			want: struct{ statusCode int }{statusCode: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockMetricStorage(ctrl)
			defer ctrl.Finish()
			handler := &UpdateMetricHandler{
				storage: mockStorage,
			}

			if tt.want.statusCode == 200 {
				mockStorage.EXPECT().UpdateMetric(tt.args.name, tt.args.value, tt.args._type).Times(1).Return(nil)
			} else if tt.want.statusCode == 501 {
				mockStorage.EXPECT().UpdateMetric(tt.args.name, tt.args.value, tt.args._type).Times(1).Return(errors.New("invalid type"))
			} else if tt.want.statusCode == 400 {
				mockStorage.EXPECT().UpdateMetric(tt.args.name, tt.args.value, tt.args._type).Times(1).Return(errors.New("invalid value"))
			}

			r := gin.Default()
			r.POST("/update/:type/:name/:value", handler.Handle)
			w := httptest.NewRecorder()
			url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", tt.args._type, tt.args.name, fmt.Sprint(tt.args.value))
			request := httptest.NewRequest(http.MethodPost, url, nil)
			r.ServeHTTP(w, request)

			assert.Equal(t, tt.want.statusCode, w.Code)
		})
	}
}