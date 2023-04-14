package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/db"
	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/commands"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger := log.New(os.Stdout, "main : ", log.LstdFlags|log.Lshortfile)
	logger.Println(`Get server's flags`)
	cmd := commands.GetServerCommands()

	storeInterval := time.NewTicker(cmd.StoreInterval)
	st := &storage.MemStorage{}

	logger.Println(`Init database`)
	db, err := db.NewDB(cmd.DataBaseDSN)
	if err != nil || cmd.DataBaseDSN == "" {
		fmt.Println("Can't compleate DB connection: ", err)
		if err != nil {
			logger.Printf("ERROR : NewDB returned the error %v\n", err)
		}
		logger.Println(`Init Memory storage`)
		st = storage.NewMemStorage(cmd.Restore, cmd.StoreInterval == 0, cmd.StoreFile)
		defer func() {
			count, err := st.SaveAllIntoFile()
			if err != nil {
				logger.Printf("ERROR : st.SaveAllIntoFile returned the error %v\n", err)
			}
			logger.Printf("saved %d metrics\n", count)
			st.ShowMetrics()
			if err = st.Close(); err != nil {
				logger.Printf("ERROR : st.Close returned the error %v\n", err)
			}
		}()
		logger.Println("Start server")
		go server.Start(st, db, cmd.Address, cmd.Key)
	} else {
		defer func() {
			if err = db.ShowMetrics(); err != nil {
				logger.Printf("ERROR : db.ShowMetrics returned the error %v\n", err)
			}
			if err = db.Close(); err != nil {
				logger.Printf("ERROR : db.Close returned the error %v\n", err)
			}
		}()
		storeInterval.Stop()
		logger.Println("Start server")
		go server.Start(db, db, cmd.Address, cmd.Key)
	}
	for {
		select {
		case <-storeInterval.C:
			count, err := st.SaveAllIntoFile()
			if err != nil {
				logger.Printf("ERROR : st.SaveAllIntoFile returned the error %v\n", err)
			}
			logger.Printf("saved %d metrics\n", count)
		case <-gracefulShutdown:
			logger.Println("Server Get Signal!")
			storeInterval.Stop()
			return
		}
	}
}
