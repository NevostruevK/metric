package metrics

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"time"
)

const MetricsCount = 29

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

func Get(cr MetricCreater) []MetricCreater {
	sM := make([]MetricCreater, 0, MetricsCount)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
/*		sM = append(sM, cr.NewGaugeMetric("Alloc", float64(mem.Alloc)))
		sM = append(sM, cr.NewGaugeMetric("BuckHashSys", float64(mem.BuckHashSys)))
		sM = append(sM, cr.NewGaugeMetric("Frees", float64(mem.Frees)))
		sM = append(sM, cr.NewGaugeMetric("GCCPUFraction", float64(mem.GCCPUFraction)))
		sM = append(sM, cr.NewGaugeMetric("GCSys", float64(mem.GCSys)))
		sM = append(sM, cr.NewGaugeMetric("HeapAlloc", float64(mem.HeapAlloc)))
		sM = append(sM, cr.NewGaugeMetric("HeapIdle", float64(mem.HeapIdle)))
		sM = append(sM, cr.NewGaugeMetric("HeapInuse", float64(mem.HeapInuse)))
		sM = append(sM, cr.NewGaugeMetric("HeapObjects", float64(mem.HeapObjects)))
		sM = append(sM, cr.NewGaugeMetric("HeapReleased", float64(mem.HeapReleased)))
		sM = append(sM, cr.NewGaugeMetric("HeapSys", float64(mem.HeapSys)))
		sM = append(sM, cr.NewGaugeMetric("LastGC", float64(mem.LastGC)))
		sM = append(sM, cr.NewGaugeMetric("Lookups", float64(mem.Lookups)))
		sM = append(sM, cr.NewGaugeMetric("MCacheInuse", float64(mem.MCacheInuse)))
		sM = append(sM, cr.NewGaugeMetric("MCacheSys", float64(mem.MCacheSys)))
		sM = append(sM, cr.NewGaugeMetric("MSpanInuse", float64(mem.MSpanInuse)))
		sM = append(sM, cr.NewGaugeMetric("MSpanSys", float64(mem.MSpanSys)))
		sM = append(sM, cr.NewGaugeMetric("Mallocs", float64(mem.Mallocs)))
		sM = append(sM, cr.NewGaugeMetric("NextGC", float64(mem.NextGC)))
		sM = append(sM, cr.NewGaugeMetric("NumForcedGC", float64(mem.NumForcedGC)))
		sM = append(sM, cr.NewGaugeMetric("NumGC", float64(mem.NumGC)))
		sM = append(sM, cr.NewGaugeMetric("OtherSys", float64(mem.OtherSys)))
		sM = append(sM, cr.NewGaugeMetric("PauseTotalNs", float64(mem.PauseTotalNs)))
		sM = append(sM, cr.NewGaugeMetric("StackInuse", float64(mem.StackInuse)))
		sM = append(sM, cr.NewGaugeMetric("StackSys", float64(mem.StackSys)))
		sM = append(sM, cr.NewGaugeMetric("Sys", float64(mem.Sys)))
*/	sM = append(sM, cr.NewGaugeMetric("TotalAlloc", float64(mem.TotalAlloc)))

	getRequestCount++
	sM = append(sM, cr.NewCounterMetric("PollCount", getRequestCount))

//		sM = append(sM, cr.NewGaugeMetric("RandomValue", getRandomFloat64()))

	return sM
}

func ResetCounter() {
	getRequestCount = 0
}
