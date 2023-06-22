package main

import (
	"context"
	"log"
	"net/http"
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

	s := &http.Server{}

	lgr.Println(`Init database`)
	db, err := db.NewDB(context.Background(), cmd.DataBaseDSN)
	if err != nil || cmd.DataBaseDSN == "" {
		lgr.Println("Can't compleate DB connection: ", err)
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
			st.ShowMetrics(context.Background())
			if err = st.Close(); err != nil {
				lgr.Printf("ERROR : st.Close returned the error %v\n", err)
			}
		}()
		s = server.NewServer(st, cmd.Address, cmd.Key)
	} else {
		defer func() {
			if err = db.ShowMetrics(context.Background()); err != nil {
				lgr.Printf("ERROR : db.ShowMetrics returned the error %v\n", err)
			}
			if err = db.Close(); err != nil {
				lgr.Printf("ERROR : db.Close returned the error %v\n", err)
			}
		}()
		storeInterval.Stop()
		s = server.NewServer(db, cmd.Address, cmd.Key)
	}
	lgr.Printf("Start server")
	go func() {
		go lgr.Println(s.ListenAndServe())
	}()
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
			if err = s.Shutdown(context.Background()); err != nil {
				lgr.Printf("ERROR : Server Shutdown error %v", err)
			} else {
				lgr.Printf("Server Shutdown ")
			}
			return
		}
	}
}
