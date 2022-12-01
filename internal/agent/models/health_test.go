package models

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/stretchr/testify/assert"
)

func TestNewGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	val := 12.23
	key := ""
	tests := []struct {
		name string
		args args
		want Metric
	}{
		{
			name: "Test create Gauge",
			args: struct {
				name  string
				value float64
			}{name: "TestGauge", value: 12.23},
			want: Metric{
				Name:       "TestGauge",
				Type:       GaugeType,
				IntValue:   nil,
				FloatValue: &val,
				Hash:       key,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGauge(tt.args.name, tt.args.value, key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	var val int64 = 56
	key := ""
	tests := []struct {
		name string
		args args
		want Metric
	}{
		{
			name: "Test create Counter",
			args: struct {
				name  string
				value int64
			}{name: "TestCounter", value: val},
			want: Metric{
				Name:       "TestCounter",
				Type:       CounterType,
				IntValue:   &val,
				FloatValue: nil,
				Hash:       key,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounter(tt.args.name, tt.args.value, key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHealth(t *testing.T) {
	type args struct {
		hashKey string
	}
	hashKey := "ff"
	tests := []struct {
		name string
		args args
		want *Health
	}{
		{
			name: "",
			args: struct{ hashKey string }{hashKey: hashKey},
			want: &Health{
				Metrics: &sync.Map{},
				Hasher:  hasher.New(hashKey),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHealth(tt.args.hashKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealth_Store(t *testing.T) {
	type fields struct {
		Metrics *sync.Map
		Hasher  *hasher.Hasher
	}
	hashKey := "ff"
	type args struct {
		id         string
		metricType MetricType
		value      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: struct {
				Metrics *sync.Map
				Hasher  *hasher.Hasher
			}{
				Metrics: &sync.Map{},
				Hasher:  hasher.New(hashKey),
			},
			args: struct {
				id         string
				metricType MetricType
				value      string
			}{id: "1", metricType: CounterType, value: "100"},
		},
		{
			name: "",
			fields: struct {
				Metrics *sync.Map
				Hasher  *hasher.Hasher
			}{
				Metrics: &sync.Map{},
				Hasher:  hasher.New(hashKey),
			},
			args: struct {
				id         string
				metricType MetricType
				value      string
			}{id: "1", metricType: GaugeType, value: "100.02"},
		},
		{
			name: "",
			fields: struct {
				Metrics *sync.Map
				Hasher  *hasher.Hasher
			}{
				Metrics: &sync.Map{},
				Hasher:  hasher.New(hashKey),
			},
			args: struct {
				id         string
				metricType MetricType
				value      string
			}{id: "1", metricType: "newType", value: "100.02"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Health{
				Metrics: tt.fields.Metrics,
				Hasher:  tt.fields.Hasher,
			}
			h.Store(tt.args.id, tt.args.metricType, tt.args.value)

			var m Metric
			if tt.args.metricType == CounterType {
				val, err := strconv.ParseInt(tt.args.value, 10, 64)
				assert.NoError(t, err)
				m = Metric{
					Name:       tt.args.id,
					Type:       tt.args.metricType,
					IntValue:   &val,
					FloatValue: nil,
					Hash:       tt.fields.Hasher.Hash(fmt.Sprintf("%s:%s:%d", tt.args.id, tt.args.metricType, val)),
				}

				mEx, ok := h.Metrics.Load(tt.args.id)
				assert.True(t, ok)
				assert.Equal(t, mEx, m)
			} else if tt.args.metricType == GaugeType {
				val, err := strconv.ParseFloat(tt.args.value, 64)
				assert.NoError(t, err)
				m = Metric{
					Name:       tt.args.id,
					Type:       tt.args.metricType,
					IntValue:   nil,
					FloatValue: &val,
					Hash:       tt.fields.Hasher.Hash(fmt.Sprintf("%s:%s:%f", tt.args.id, tt.args.metricType, val)),
				}

				mEx, ok := h.Metrics.Load(tt.args.id)
				assert.True(t, ok)
				assert.Equal(t, mEx, m)
			} else {
				_, ok := h.Metrics.Load(tt.args.id)
				assert.False(t, ok)
			}
		})
	}
}
