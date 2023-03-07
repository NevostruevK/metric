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
        inMetrics := make([]metrics.Metric, 0, metrics.MetricsCount*(reportInterval/pollInterval+1))
        outMetrics := make([]metrics.Metric, metrics.MetricsCount*(reportInterval/pollInterval+1))
        var size int
        for {
                select {
                case <-pollTicker.C:
                        inMetrics = append(inMetrics, metrics.Get()...)
                case <-reportTicker.C:
                        size = copy(outMetrics, inMetrics)
                        inMetrics = nil
                        client.SendMetrics(outMetrics, size)
                case <-gracefulShutdown:
                        pollTicker.Stop()
                        reportTicker.Stop()
                        fmt.Println("Get Signal")
                        return
                }
        }

}

