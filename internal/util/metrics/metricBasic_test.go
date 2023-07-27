package metrics

import (
	"reflect"
	"testing"
)

func TestBasicMetric_ConvertToMetrics(t *testing.T) {
	//	bm := NewCounterMetric("testCounter",12345)
	tests := []struct {
		name string
		m    *BasicMetric
		want Metrics
	}{
		{
			name: "normal counter test",
			m:    NewCounterMetric("testCounter", 12345),
			want: NewJSONCounterMetric("testCounter", 12345),
		},
		{
			name: "normal gauge test",
			m:    NewGaugeMetric("testGauge", 1.2345),
			want: NewJSONGaugeMetric("testGauge", 1.2345),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ConvertToMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicMetric.ConvertToMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
