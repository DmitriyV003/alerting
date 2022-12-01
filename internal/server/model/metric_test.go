package model

import (
	"reflect"
	"testing"
)

func TestNewCounter(t *testing.T) {
	var val int64 = 120
	type args struct {
		name string
		val  int64
	}
	tests := []struct {
		name string
		args args
		want *Metric
	}{
		{
			name: "",
			args: struct {
				name string
				val  int64
			}{
				name: "fg",
				val:  val,
			},
			want: &Metric{
				Name:       "fg",
				Type:       CounterType,
				IntValue:   &val,
				FloatValue: nil,
				Hash:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounter(tt.args.name, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGauge(t *testing.T) {
	var val = 120.045
	type args struct {
		name string
		val  float64
	}
	tests := []struct {
		name string
		args args
		want *Metric
	}{
		{
			name: "",
			args: struct {
				name string
				val  float64
			}{
				name: "fg",
				val:  val,
			},
			want: &Metric{
				Name:       "fg",
				Type:       GaugeType,
				IntValue:   nil,
				FloatValue: &val,
				Hash:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGauge(tt.args.name, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMetric(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Metric
	}{
		{
			name: "",
			args: struct{ name string }{name: "fg"},
			want: &Metric{
				Name: "fg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMetric(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
