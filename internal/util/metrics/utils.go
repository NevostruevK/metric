package metrics

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

const (
	MetricsCount      = 29
	ExtraMetricsCount = 3
)

type gauge float64
type counter int64

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func roundGauge(f float64) string {
	s := fmt.Sprintf("%.3f", f)
	r := regexp.MustCompile("0{1,2}$")
	return r.ReplaceAllString(s, "")
}

func IsMetricType(checkType string) bool {
	if checkType != Gauge && checkType != Counter {
		return false
	}
	return true
}

type MetricCreater interface {
	NewGaugeMetric(name string, value float64) MetricCreater
	NewCounterMetric(name string, value int64) MetricCreater
}

var getRequestCount int64 = 0

func getRandomFloat64() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

// func GetAdvanced() ([]MetricCreater, error){
func GetAdvanced() ([]Metrics, error) {

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	sM := make([]Metrics, 0, ExtraMetricsCount)
	sM = append(sM, NewJSONGaugeMetric("TotalMemory", float64(v.Total)))
	sM = append(sM, NewJSONGaugeMetric("FreeMemory", float64(v.Free)))
	sM = append(sM, NewJSONGaugeMetric("CPUutilization1", v.UsedPercent))
	return sM, nil
}

// func Get(cr MetricCreater) []MetricCreater {
func Get() []Metrics {
	sM := make([]Metrics, 0, MetricsCount)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	sM = append(sM, NewJSONGaugeMetric("Alloc", float64(mem.Alloc)))
	sM = append(sM, NewJSONGaugeMetric("BuckHashSys", float64(mem.BuckHashSys)))
	sM = append(sM, NewJSONGaugeMetric("Frees", float64(mem.Frees)))
	sM = append(sM, NewJSONGaugeMetric("GCCPUFraction", float64(mem.GCCPUFraction)))
	sM = append(sM, NewJSONGaugeMetric("GCSys", float64(mem.GCSys)))
	sM = append(sM, NewJSONGaugeMetric("HeapAlloc", float64(mem.HeapAlloc)))
	sM = append(sM, NewJSONGaugeMetric("HeapIdle", float64(mem.HeapIdle)))
	sM = append(sM, NewJSONGaugeMetric("HeapInuse", float64(mem.HeapInuse)))
	sM = append(sM, NewJSONGaugeMetric("HeapObjects", float64(mem.HeapObjects)))
	sM = append(sM, NewJSONGaugeMetric("HeapReleased", float64(mem.HeapReleased)))
	sM = append(sM, NewJSONGaugeMetric("HeapSys", float64(mem.HeapSys)))
	sM = append(sM, NewJSONGaugeMetric("LastGC", float64(mem.LastGC)))
	sM = append(sM, NewJSONGaugeMetric("Lookups", float64(mem.Lookups)))
	sM = append(sM, NewJSONGaugeMetric("MCacheInuse", float64(mem.MCacheInuse)))
	sM = append(sM, NewJSONGaugeMetric("MCacheSys", float64(mem.MCacheSys)))
	sM = append(sM, NewJSONGaugeMetric("MSpanInuse", float64(mem.MSpanInuse)))
	sM = append(sM, NewJSONGaugeMetric("MSpanSys", float64(mem.MSpanSys)))
	sM = append(sM, NewJSONGaugeMetric("Mallocs", float64(mem.Mallocs)))
	sM = append(sM, NewJSONGaugeMetric("NextGC", float64(mem.NextGC)))
	sM = append(sM, NewJSONGaugeMetric("NumForcedGC", float64(mem.NumForcedGC)))
	sM = append(sM, NewJSONGaugeMetric("NumGC", float64(mem.NumGC)))
	sM = append(sM, NewJSONGaugeMetric("OtherSys", float64(mem.OtherSys)))
	sM = append(sM, NewJSONGaugeMetric("PauseTotalNs", float64(mem.PauseTotalNs)))
	sM = append(sM, NewJSONGaugeMetric("StackInuse", float64(mem.StackInuse)))
	sM = append(sM, NewJSONGaugeMetric("StackSys", float64(mem.StackSys)))
	sM = append(sM, NewJSONGaugeMetric("Sys", float64(mem.Sys)))
	sM = append(sM, NewJSONGaugeMetric("TotalAlloc", float64(mem.TotalAlloc)))

	getRequestCount++
	sM = append(sM, NewJSONCounterMetric("PollCount", getRequestCount))

	sM = append(sM, NewJSONGaugeMetric("RandomValue", getRandomFloat64()))

	return sM
}

func ResetCounter() {
	getRequestCount = 0
}
