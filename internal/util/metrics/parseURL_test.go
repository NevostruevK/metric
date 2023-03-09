package metrics_test

import (
	"reflect"
	"testing"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

func TestURLToMetric(t *testing.T) {
	tests := []struct {
		name    string
		url    	string
		want    *metrics.Metric
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ok simple test",
			url: "/update/gauge/RandomValue/0.949516",
			want: metrics.NewGaugeMetric("RandomValue", 0.949516),
			wantErr: false,
		},
		{
			name: "err type field",
			url: "/update/gaug/RandomValue/0.949516",
			want: nil,
			wantErr: true,
		},
		{
			name: "err type value",
			url: "/update/counter/RandomValue/0.949516",
			want: nil,
			wantErr: true,
		},
		{
			name: "err too many fields",
			url: "/update/gauge/RandomValue/0.949516/",
			want: nil,
			wantErr: true,
		},
		{
			name: "err too few fields",
			url: "update/gauge/RandomValue/0.949516/",
			want: nil,
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := metrics.URLToMetric(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLToMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URLToMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
