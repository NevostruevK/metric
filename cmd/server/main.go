package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/commands"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cmd := commands.GetServerCommands()
	fmt.Printf("Server get command %+v\n", cmd)

	//	server.SetAddress(cmd.Address)
	st := storage.NewMemStorage(cmd.Restore, cmd.StoreInterval == 0, cmd.StoreFile)
	storeInterval := time.NewTicker(cmd.StoreInterval)

	go server.Start(st, cmd.Address, cmd.Key)

	for {
		select {
		case <-storeInterval.C:
			count, err := st.SaveAllIntoFile()
			fmt.Printf("Saved to file %d metrics, error %v\n", count, err)
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
