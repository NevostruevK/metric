package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/storage"
)

func main() {
	
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	storage := storage.NewMemStorage()
/*	storage := &server.SaveStorage{
		Data: make(map[string]metrics.Metric),
	}
*/
	go server.Start(storage)
	
	<-gracefulShutdown
	fmt.Println("Server Get Signal!")
	storage.ShowMetrics()
}
