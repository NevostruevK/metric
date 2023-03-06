package main

import (
	"time"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

//const SERVER_ADDRES = "http://localhost:8080/"
//const SERVER_ADDRES = "http://127.0.0.1:8080/"
const pollInterval = 2
const reportInterval = 10

func main() {
	pollTicker := time.NewTicker(pollInterval*time.Second)
	reportTicker := time.NewTicker(reportInterval*time.Second)
	inMetrics := make([]metrics.Metric,0,metrics.MetricsCount * (reportInterval/pollInterval+1))
	outMetrics := make([]metrics.Metric,metrics.MetricsCount * (reportInterval/pollInterval+1))
	var size int
//	var outMetrics []metrics.Metric = nil
//	metrics.Get()
	for{
	select{
	case <- pollTicker.C :
		inMetrics = append(inMetrics, metrics.Get()...)
	case <- reportTicker.C:
//		outMetrics = inMetrics
		size = copy(outMetrics, inMetrics)
//		fmt.Println("=================copy================= :", size)
		inMetrics = nil
		client.SendMetrics(outMetrics, size)
	}
}
//	sM = metrics.Get()
//	sM = append(sM, metrics.Get()...)
//	fmt.Println("client.SendMetric(sM)")
//	client.SendMetric(sM[0])
//	client.SendMetrics(sM)
//	client.HelloFromClient()

}

