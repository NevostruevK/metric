package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/caarlos0/env/v7"
)

type environment struct{
	Address 		string 			`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile		string  		`env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	StoreInterval	time.Duration  	`env:"STORE_INTERVAL" envDefault:"300s"`
	Restore			bool			`env:"RESTORE" envDefault:"true"`
}

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	en := environment{}
	if err := env.Parse(&en); err != nil{
		fmt.Printf("Server read environment with the error: %+v\n", err)
	}
	server.SetAddress(en.Address)
	st := storage.NewMemStorage(en.Restore, en.StoreInterval==0, en.StoreFile)
	storeInterval := time.NewTicker(en.StoreInterval)

	go server.Start(st)

	for {
		select {
		case <-storeInterval.C:
			count, err := st.SaveAllIntoFile()
			fmt.Printf("Saved to file %d metrics, error %v\n",count,err)
		case <-gracefulShutdown:
			fmt.Println("Server Get Signal!")
			storeInterval.Stop()
			st.SaveAllIntoFile()
			st.Close()
			st.ShowMetrics()
			return
		}
	}
}
