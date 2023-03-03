package client

import (
	"testing"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

func TestSendMetric(t *testing.T) {
	type args struct {
		sM metrics.Metric
//sM metrics.Metric
//		sM metrics.MMM
}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple ok",
//			args : args{metrics.MMM{Name: "Alloc", MmType: "gauge", MgValue: 12345, McValue: 0}},	
			args : args{metrics.Metric{"Alloc", "gauge",  12345, 0},},	
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			SendMetric(tt.args.sM)
		})
	}
}
