package metrics

import (
	"fmt"
	"runtime"
)

type gauge float64
type counter int64

type MetricType string
type MetricName string
const (
	GAUGE  = MetricType("gauge")
	COUNTER = MetricType("counter")
)
type Metric struct{
	Name string
	Type MetricType
	MgValue gauge
	McValue counter
}
type MMM struct{
	Mname string
	MTypeM MetricType
	MgValue gauge
	MMcValue counter
}
func (m Metric) GetType() string{
	return string(m.Type)
}
func (m Metric) GetName() string{
	return string(m.Name)
}
func (m Metric) GetValue() string{
	if (m.Type==COUNTER){
		return fmt.Sprintf("%d", m.McValue)
	}
	return fmt.Sprintf("%f", m.MgValue)
//	string(m.GetValue())
}
const METRIX_COUNT = 20
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

func Get() []Metric{
	sM := make([]Metric, 0, METRIX_COUNT) 
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)    
	sM = append(sM, Metric{"Alloc",GAUGE, gauge(mem.Alloc), 0})
	sM = append(sM, Metric{"BuckHashSys",GAUGE, gauge(mem.BuckHashSys), 0})
	sM = append(sM, Metric{"Frees",GAUGE, gauge(mem.Frees), 0})
	sM = append(sM, Metric{"GCCPUFraction",GAUGE, gauge(mem.GCCPUFraction), 0})
/*	fmt.Println("mem.Alloc:", mem.Alloc)    
	fmt.Println("mem.TotalAlloc:", mem.TotalAlloc)    
	fmt.Println("mem.HeapAlloc:", mem.HeapAlloc)    
	fmt.Println("mem.NumGC:", mem.NumGC)	
*/
	for i, m:= range sM {
		fmt.Println(i, ":", m.Name, m.Type, m.MgValue, m.McValue)
	}
	return sM
}