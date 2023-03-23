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

//const pollInterval = 2
//const reportInterval = 10

type environment struct{
	Address 		string 			`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval	time.Duration	`env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval	time.Duration	`env:"POLL_INTERVAL" envDefault:"2s"`
}

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	en := environment{}
	if err := env.Parse(&en); err!=nil{
		fmt.Printf("Agent read environment with the error: %+v\n", err)
	}
	fmt.Printf("Get environment %+v\n",en)
	client.SetAddress(en.Address)
	pollTicker := time.NewTicker(en.PollInterval)
	reportTicker := time.NewTicker(en.ReportInterval)

	sM := make([]metrics.MetricCreater, 0, metrics.MetricsCount*(en.ReportInterval/en.PollInterval+2))
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
