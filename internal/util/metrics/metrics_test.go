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
			m:    *metrics.Float64ToGauge(1.23456789).NewMetric("GaugeMetric"),
			//                      want: "gauge/GaugeMetric/1.23",
			want: fmt.Sprintf("gauge/GaugeMetric/%f", 1.23456789),
		},
		{
			name: "simple counter metric",
			m:    *metrics.Int64ToGauge(23456789).NewMetric("CounterMetric"),
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

func TestMetric_AddMetricValue(t *testing.T) {
	tests := []struct {
		name    string
		m       metrics.Metric
		new     metrics.Metric
		want    metrics.Metric
		wantErr bool
	}{
		{
                        name:    "simple ok counter",
                        m:       *metrics.NewCounterMetric("okCounter 1+5",1), 
                        new:     *metrics.NewCounterMetric("add 5",5), 
                        want:    *metrics.NewCounterMetric("okCounter 1+5",6), 
                        wantErr: false,
                 },
                 {
                        name:    "simple err different types",
                        m:       *metrics.NewCounterMetric("errCounter 1+5",1), 
                        new:     *metrics.NewGaugeMetric("add 5",5), 
                        want:    *metrics.NewCounterMetric("errCounter 1+5",1), 
                        wantErr: true,
                 },
  
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.AddMetricValue(tt.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.AddMetricValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Metric.AddMetricValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
