package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dmitriy/alerting/internal/server/mocks"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/transformers"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetricByTypeAndNameHandler_Handle(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			name string
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test successful update counter",
			args: struct {
				name string
			}{
				name: "testCounter",
			},
			want: struct {
				statusCode int
			}{statusCode: 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockMetricStorage(ctrl)
			defer ctrl.Finish()
			handler := &GetMetricByTypeAndNameHandler{
				storage:           mockStorage,
				metricTransformer: transformers.NewTransformer("ddd"),
			}

			router.Post("/value", handler.Handle)

			mockStorage.EXPECT().GetByNameAndType(context.Background(), tt.args.name, "gauge").Times(1).Return(model.NewGauge(tt.args.name, 45.2), nil)

			body := metricRequest{
				Name: tt.args.name,
				Type: "gauge",
			}
			js, err := json.Marshal(body)
			assert.NoError(t, err)

			newRequest, err := http.NewRequest(http.MethodPost, "http://localhost:8080/value", bytes.NewBuffer(js))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, newRequest)
			res := rr.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
