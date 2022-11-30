package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
			name: "test successful update counter",
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
			name: "test successful update gauge",
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
			router := chi.NewRouter()
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockMetricStorage(ctrl)
			defer ctrl.Finish()
			handler := &UpdateMetricHandler{
				storage: mockStorage,
			}

			router.Post("/update/{type}/{name}/{value}", handler.Handle)

			if tt.want.statusCode == 200 {
				mockStorage.EXPECT().UpdateOrCreate(context.Background(), tt.args.name, tt.args.value, tt.args._type).Times(1).Return(nil)
			} else if tt.want.statusCode == 501 {
				mockStorage.EXPECT().UpdateOrCreate(context.Background(), tt.args.name, tt.args.value, tt.args._type).Times(1).Return(applicationerrors.ErrInvalidType)
			} else if tt.want.statusCode == 400 {
				mockStorage.EXPECT().UpdateOrCreate(context.Background(), tt.args.name, tt.args.value, tt.args._type).Times(1).Return(applicationerrors.ErrInvalidValue)
			}

			url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", tt.args._type, tt.args.name, fmt.Sprint(tt.args.value))
			newRequest, err := http.NewRequest(http.MethodPost, url, nil)
			rr := httptest.NewRecorder()

			if err != nil {
				return
			}

			router.ServeHTTP(rr, newRequest)
			res := rr.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
