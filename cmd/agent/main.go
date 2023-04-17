package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	lgr.Println(`Get server's flags`)

	cmd, err := commands.GetAgentCommands()
	logger.LogCommands(cmd, false, err)

	a := client.NewAgent(cmd.Address, cmd.Key)
	pollTicker := time.NewTicker(cmd.PollInterval)
	reportTicker := time.NewTicker(cmd.ReportInterval)

	sM := make([]metrics.MetricCreater, 0, metrics.MetricsCount*(cmd.ReportInterval/cmd.PollInterval+2))
	mInit := metrics.Metrics{}

	for {
		select {
		case <-pollTicker.C:
			lgr.Println("Get Metrics")
			sM = append(sM, metrics.Get(&mInit)...)
		case <-reportTicker.C:
			lgr.Println("Send Metric: ", len(sM))
			sendCount := client.SendMetrics(a, sM)
			if sendCount == len(sM) {
				metrics.ResetCounter()
				sM = nil
				break
			}
			lgr.Println("Sent ", sendCount, "metrics from ", len(sM))
			sM = sM[sendCount:]
		case <-gracefulShutdown:
			lgr.Println("Get Agent Signal!")
			pollTicker.Stop()
			reportTicker.Stop()
			return
		}
	}
}
