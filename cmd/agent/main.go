package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

const pollInterval = 2
const reportInterval = 10

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	pollTicker := time.NewTicker(pollInterval * time.Second)
	reportTicker := time.NewTicker(reportInterval * time.Second)

	sM := make([]metrics.MetricCreater, 0, metrics.MetricsCount*(reportInterval/pollInterval+1))
//	mInit := metrics.Metric{}
	mInit := metrics.Metrics{}
		
	for {
		select {
		case <-pollTicker.C:
			fmt.Println("Get Metric")
			sM = append(sM, metrics.Get(&mInit)...)
		case <-reportTicker.C:
			fmt.Println("Send Metric")
			client.SendMetrics(sM)
			sM = nil
		case <-gracefulShutdown:
			pollTicker.Stop()
			reportTicker.Stop()
			fmt.Println("Get Agent Signal!")
			return
		}
	}
}
