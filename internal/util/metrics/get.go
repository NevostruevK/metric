package metrics

import (
	"math/rand"
	"runtime"
	"time"
)
var getRequestCount int64 = 0
func getRandomValue() Metric{
	rand.Seed(time.Now().UnixNano())
	f:= rand.Float64()
	return gauge(f).NewMetric("RandomValue")
}

func Get() []Metric{
	sM := make([]Metric, 0, MetricsCount) 
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)    
	sM = append(sM, gauge(mem.Alloc).NewMetric("Alloc"))
	sM = append(sM, gauge(mem.BuckHashSys).NewMetric("BuckHashSys"))
	sM = append(sM, gauge(mem.Frees).NewMetric("Frees"))
	sM = append(sM, gauge(mem.GCCPUFraction).NewMetric("GCCPUFraction"))
	sM = append(sM, gauge(mem.GCSys).NewMetric("GCSys"))
	sM = append(sM, gauge(mem.HeapAlloc).NewMetric("HeapAlloc"))
	sM = append(sM, gauge(mem.HeapIdle).NewMetric("HeapIdle"))
	sM = append(sM, gauge(mem.HeapInuse).NewMetric("HeapInuse"))
	sM = append(sM, gauge(mem.HeapObjects).NewMetric("HeapObjects"))
	sM = append(sM, gauge(mem.HeapReleased).NewMetric("HeapReleased"))
	sM = append(sM, gauge(mem.HeapSys).NewMetric("HeapSys"))
	sM = append(sM, gauge(mem.LastGC).NewMetric("LastGC"))
	sM = append(sM, gauge(mem.Lookups).NewMetric("Lookups"))
	sM = append(sM, gauge(mem.MCacheInuse).NewMetric("MCacheInuse"))
	sM = append(sM, gauge(mem.MCacheSys).NewMetric("MCacheSys"))
	sM = append(sM, gauge(mem.MSpanInuse).NewMetric("MSpanInuse"))
	sM = append(sM, gauge(mem.MSpanSys).NewMetric("MSpanSys"))
	sM = append(sM, gauge(mem.Mallocs).NewMetric("Mallocs"))
	sM = append(sM, gauge(mem.NextGC).NewMetric("NextGC"))
	sM = append(sM, gauge(mem.NumForcedGC).NewMetric("NumForcedGC"))
	sM = append(sM, gauge(mem.NumGC).NewMetric("NumGC"))
	sM = append(sM, gauge(mem.OtherSys).NewMetric("OtherSys"))
	sM = append(sM, gauge(mem.PauseTotalNs).NewMetric("PauseTotalNs"))
	sM = append(sM, gauge(mem.StackInuse).NewMetric("StackInuse"))
	sM = append(sM, gauge(mem.StackSys).NewMetric("StackSys"))
	sM = append(sM, gauge(mem.Sys).NewMetric("Sys"))
	sM = append(sM, gauge(mem.TotalAlloc).NewMetric("TotalAlloc"))

	sM = append(sM, getRandomValue())

	sM = append(sM, counter(getRequestCount).NewMetric("PollCount"))
	getRequestCount++

/*	for i, m:= range sM {
		fmt.Println(i, ":",m.String())
	}
	fmt.Println("-----------------")
*/	return sM
}
/*
const (
	Alloc			MetricName("Alloc")
	BuckHashSys 	MetricName("BuckHashSys")
	Frees			MetricName("Frees")
)*/
/*
Имя метрики: "Alloc", тип: gauge
Имя метрики: "BuckHashSys", тип: gauge
Имя метрики: "Frees", тип: gauge
Имя метрики: "GCCPUFraction", тип: gauge
Имя метрики: "GCSys", тип: gauge
Имя метрики: "HeapAlloc", тип: gauge
Имя метрики: "HeapIdle", тип: gauge
Имя метрики: "HeapInuse", тип: gauge
Имя метрики: "HeapObjects", тип: gauge
Имя метрики: "HeapReleased", тип: gauge
Имя метрики: "HeapSys", тип: gauge
Имя метрики: "LastGC", тип: gauge
Имя метрики: "Lookups", тип: gauge
Имя метрики: "MCacheInuse", тип: gauge
Имя метрики: "MCacheSys", тип: gauge
Имя метрики: "MSpanInuse", тип: gauge
Имя метрики: "MSpanSys", тип: gauge
Имя метрики: "Mallocs", тип: gauge
Имя метрики: "NextGC", тип: gauge
Имя метрики: "NumForcedGC", тип: gauge
Имя метрики: "NumGC", тип: gauge
Имя метрики: "OtherSys", тип: gauge
Имя метрики: "PauseTotalNs", тип: gauge
Имя метрики: "StackInuse", тип: gauge
Имя метрики: "StackSys", тип: gauge
Имя метрики: "Sys", тип: gauge
Имя метрики: "TotalAlloc", тип: gauge
*/

