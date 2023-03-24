package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	cmd := commands.GetAgentCommands()

	fmt.Printf("Agent get command %+v\n",cmd)
	client.SetAddress(cmd.Address)
	pollTicker := time.NewTicker(cmd.PollInterval)
	reportTicker := time.NewTicker(cmd.ReportInterval)

	sM := make([]metrics.MetricCreater, 0, metrics.MetricsCount*(cmd.ReportInterval/cmd.PollInterval+2))
//	mInit := metrics.Metric{}
	mInit := metrics.Metrics{}
		
	for {
		select {
		case <-pollTicker.C:
			fmt.Println("Get Metric")
			sM = append(sM, metrics.Get(&mInit)...)
		case <-reportTicker.C:
			fmt.Println("Send Metric: ",len(sM))
			sendCount := client.SendMetrics(sM)
			if sendCount == len(sM){
				metrics.ResetCounter()
				sM = nil
				break
			}
			fmt.Println("Send not All metrics ",sendCount," from ", len(sM))
			sM = sM[sendCount:]
		case <-gracefulShutdown:
			fmt.Println("Get Agent Signal!")
			pollTicker.Stop()
			reportTicker.Stop()
			return
		}
	}
}
