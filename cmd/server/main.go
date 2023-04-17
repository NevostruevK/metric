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
	"github.com/NevostruevK/metric/internal/util/logger"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	lgr.Println(`Get server's flags`)

	cmd, err := commands.GetServerCommands()
	logger.LogCommands(cmd, true, err)


	storeInterval := time.NewTicker(cmd.StoreInterval)
	st := &storage.MemStorage{}

	lgr.Println(`Init database`)
	db, err := db.NewDB(cmd.DataBaseDSN)
	if err != nil || cmd.DataBaseDSN == "" {
		fmt.Println("Can't compleate DB connection: ", err)
		if err != nil {
			lgr.Printf("ERROR : NewDB returned the error %v\n", err)
		}
		lgr.Println(`Init Memory storage`)
		st = storage.NewMemStorage(cmd.Restore, cmd.StoreInterval == 0, cmd.StoreFile)
		defer func() {
			count, err := st.SaveAllIntoFile()
			if err != nil {
				lgr.Printf("ERROR : st.SaveAllIntoFile returned the error %v\n", err)
			}
			lgr.Printf("saved %d metrics\n", count)
			st.ShowMetrics()
			if err = st.Close(); err != nil {
				lgr.Printf("ERROR : st.Close returned the error %v\n", err)
			}
		}()
		lgr.Println("Start server")
		go server.Start(st, db, cmd.Address, cmd.Key)
	} else {
		defer func() {
			if err = db.ShowMetrics(); err != nil {
				lgr.Printf("ERROR : db.ShowMetrics returned the error %v\n", err)
			}
			if err = db.Close(); err != nil {
				lgr.Printf("ERROR : db.Close returned the error %v\n", err)
			}
		}()
		storeInterval.Stop()
		lgr.Println("Start server")
		go server.Start(db, db, cmd.Address, cmd.Key)
	}
	for {
		select {
		case <-storeInterval.C:
			count, err := st.SaveAllIntoFile()
			if err != nil {
				lgr.Printf("ERROR : st.SaveAllIntoFile returned the error %v\n", err)
			}
			lgr.Printf("saved %d metrics\n", count)
		case <-gracefulShutdown:
			lgr.Println("Server Get Signal!")
			storeInterval.Stop()
			return
		}
	}
}
