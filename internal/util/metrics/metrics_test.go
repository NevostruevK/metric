package metrics_test

import (
	"fmt"
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