package metrics

import (
	"math/rand"
	"runtime"
	"time"
)

var getRequestCount int64 = 0

func getRandomValue() *Metric {
	rand.Seed(time.Now().UnixNano())
	f := rand.Float64()
	return gauge(f).NewMetric("RandomValue")
}

func Get() []Metric {
	sM := make([]Metric, 0, MetricsCount)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	sM = append(sM, *gauge(mem.Alloc).NewMetric("Alloc"))
	sM = append(sM, *gauge(mem.BuckHashSys).NewMetric("BuckHashSys"))
	sM = append(sM, *gauge(mem.Frees).NewMetric("Frees"))
	sM = append(sM, *gauge(mem.GCCPUFraction).NewMetric("GCCPUFraction"))
	sM = append(sM, *gauge(mem.GCSys).NewMetric("GCSys"))
	sM = append(sM, *gauge(mem.HeapAlloc).NewMetric("HeapAlloc"))
	sM = append(sM, *gauge(mem.HeapIdle).NewMetric("HeapIdle"))
	sM = append(sM, *gauge(mem.HeapInuse).NewMetric("HeapInuse"))
	sM = append(sM, *gauge(mem.HeapObjects).NewMetric("HeapObjects"))
	sM = append(sM, *gauge(mem.HeapReleased).NewMetric("HeapReleased"))
	sM = append(sM, *gauge(mem.HeapSys).NewMetric("HeapSys"))
	sM = append(sM, *gauge(mem.LastGC).NewMetric("LastGC"))
	sM = append(sM, *gauge(mem.Lookups).NewMetric("Lookups"))
	sM = append(sM, *gauge(mem.MCacheInuse).NewMetric("MCacheInuse"))
	sM = append(sM, *gauge(mem.MCacheSys).NewMetric("MCacheSys"))
	sM = append(sM, *gauge(mem.MSpanInuse).NewMetric("MSpanInuse"))
	sM = append(sM, *gauge(mem.MSpanSys).NewMetric("MSpanSys"))
	sM = append(sM, *gauge(mem.Mallocs).NewMetric("Mallocs"))
	sM = append(sM, *gauge(mem.NextGC).NewMetric("NextGC"))
	sM = append(sM, *gauge(mem.NumForcedGC).NewMetric("NumForcedGC"))
	sM = append(sM, *gauge(mem.NumGC).NewMetric("NumGC"))
	sM = append(sM, *gauge(mem.OtherSys).NewMetric("OtherSys"))
	sM = append(sM, *gauge(mem.PauseTotalNs).NewMetric("PauseTotalNs"))
	sM = append(sM, *gauge(mem.StackInuse).NewMetric("StackInuse"))
	sM = append(sM, *gauge(mem.StackSys).NewMetric("StackSys"))
	sM = append(sM, *gauge(mem.Sys).NewMetric("Sys"))
	sM = append(sM, *gauge(mem.TotalAlloc).NewMetric("TotalAlloc"))

	sM = append(sM, *counter(getRequestCount).NewMetric("PollCount"))
	getRequestCount++

	sM = append(sM, *getRandomValue())

	return sM
}
