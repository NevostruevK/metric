package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/caarlos0/env/v7"
)

const pollInterval = 2
const reportInterval = 10

type environment struct{
	Address 		string 			`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval	int				`env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval	int				`env:"POLL_INTERVAL" envDefault:"2"`
}

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	en := environment{}
	if err := env.Parse(&en); err!=nil{
		fmt.Printf("Agent read environment with the error: %+v\n", err)
	}
	client.SetAddress(en.Address)
	pollTicker := time.NewTicker(time.Duration(en.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(en.ReportInterval) * time.Second)

	sM := make([]metrics.MetricCreater, 0, metrics.MetricsCount*(reportInterval/pollInterval+2))
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
			pollTicker.Stop()
			reportTicker.Stop()
			fmt.Println("Get Agent Signal!")
			return
		}
	}
}
