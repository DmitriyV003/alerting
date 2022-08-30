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
