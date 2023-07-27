package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func Test_IsMetricType(t *testing.T) {
	type args struct {
		typ string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test Countrer",
			args: args{Counter},
			want: true,
		},
		{
			name: "test Gauge",
			args: args{Gauge},
			want: true,
		},
		{
			name: "test empty",
			args: args{""},
			want: false,
		},
		{
			name: "test no metric type",
			args: args{"Integer"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMetricType(tt.args.typ); got != tt.want {
				t.Errorf("IsMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	got := Get()
	require.NotNil(t, got)
	assert.True(t, len(got) == MetricsCount)
	for _, m := range got {
		assert.True(t, IsMetricType(m.MType))
		assert.NotEmpty(t, m.Name())
		assert.True(t, m.Delta != nil || m.Value != nil)
	}
}

func TestGetAdvanced(t *testing.T) {
	got, err := GetAdvanced()
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, len(got) == ExtraMetricsCount)
	for _, m := range got {
		assert.True(t, IsMetricType(m.MType))
		assert.NotEmpty(t, m.Name())
		assert.True(t, m.Delta != nil || m.Value != nil)
	}
}

func TestResetCounter(t *testing.T) {
	getRequestCount = 12345
	ResetCounter()
	assert.True(t, getRequestCount == 0)
}

func TestRandomFloat64(t *testing.T) {
	random1 := getRandomFloat64()
	time.Sleep(time.Millisecond)
	random2 := getRandomFloat64()
	assert.NotEqual(t, random1, random2)
}
