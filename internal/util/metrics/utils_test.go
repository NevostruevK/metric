package metrics

import "testing"

func Test_roundGauge(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 12.34567890",
			args: args{f: 12.34567890},
			want: "12.346",
		},
		{
			name: "test 12.340",
			args: args{f: 12.340},
			want: "12.34",
		},
		{
			name: "test 12.000",
			args: args{f: 12.000},
			want: "12.0",
		},
		{
			name: "test 12.103",
			args: args{f: 12.103},
			want: "12.103",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundGauge(tt.args.f); got != tt.want {
				t.Errorf("roundGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}
