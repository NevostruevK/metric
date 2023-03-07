package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

func main() {
	m, err := metrics.URLToMetric("/update/gauge/name123/0.123")
	if err!=nil{
		fmt.Println("convert URL to metric error",err)
		return
	}
	fmt.Println(m.String())
	os.Exit(1)
	
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	storage := &server.SaveStorage{
		Data: make(map[string]metrics.Metric),
	}

	go server.Start(storage)
	
	<-gracefulShutdown
	fmt.Println("Server Get Signal!")
	for i, m := range storage.Data{
		fmt.Println(i, " : ",m.String())
	}
}
