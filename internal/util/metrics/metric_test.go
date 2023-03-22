package metrics_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

func TestMetric_String(t *testing.T) {
	tests := []struct {
		name string
		m    metrics.Metric
		want string
	}{
		{
			name: "simple gauge metric",
			m:    *metrics.NewGaugeMetric("GaugeMetric", 1.23456789),
			want: fmt.Sprintf("gauge/GaugeMetric/%.3f", 1.23456789),
		},
		{
			name: "simple counter metric",
			m:    *metrics.NewCounterMetric("CounterMetric", 23456789),
			want: "counter/CounterMetric/23456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Metric.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_AddCounterValue(t *testing.T) {
	tests := []struct {
		name    string
		m       *metrics.Metric
		value   int64
		want    *metrics.Metric
		wantErr bool
	}{
		{
			name:    "simple ok counter",
			m:       metrics.NewCounterMetric("okCounter 1+5", 1),
			value:   5,
			want:    metrics.NewCounterMetric("okCounter 1+5", 6),
			wantErr: false,
		},
		{
			name:    "simple err gauge type",
			m:       metrics.NewGaugeMetric("errGauge 1+5", 1),
			value:   5,
			want:    metrics.NewGaugeMetric("errGauge 1+5", 1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.AddCounterValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.AddMetricValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Metric.AddMetricValue() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}
