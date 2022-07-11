package models

import (
	"reflect"
	"testing"
)

func TestNewGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args args
		want Gauge
	}{
		{
			name: "Test create Gauge",
			args: struct {
				name  string
				value float64
			}{name: "TestGauge", value: 12.23},
			want: struct {
				Metric
				Value float64
			}{Metric: Metric{Name: "TestGauge"}, Value: 12.23},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGauge(tt.args.name, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}
