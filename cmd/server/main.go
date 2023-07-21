package main

import (
	"context"
	"fmt"
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

var (
	buildVersion = "N/A"
	buildData    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx := context.Background()
	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	/*
		cmd1 := commands.NewServerOptions()
		cmd1.ReadConfig("config.json")
		fmt.Println(cmd1)
	*/
	fmt.Println("---------------")
	lgr.Println("Build version : " + buildVersion)
	lgr.Println("Build data    : " + buildData)
	lgr.Println("Build commit  : " + buildCommit)

	lgr.Println(`Get server's flags`)

	cfg := commands.GetServerConfig()
	//	cmd, err := commands.GetServerCommands()
	logger.LogCommands(cfg, true)

	storeInterval := time.NewTicker(cfg.StoreInterval.Duration)
	st := &storage.MemStorage{}

	s := &http.Server{}

	lgr.Println(`Init database`)
	db, err := db.NewDB(ctx, cfg.DataBaseDSN)
	if err != nil || cfg.DataBaseDSN == "" {
		lgr.Println("Can't compleate DB connection: ", err)
		if err != nil {
			lgr.Printf("ERROR : NewDB returned the error %v\n", err)
		}
		lgr.Println(`Init Memory storage`)
		st = storage.NewMemStorage(cfg.Restore, cfg.StoreInterval.Duration == 0, cfg.StoreFile)
		defer func() {
			count, errSave := st.SaveAllIntoFile()
			if errSave != nil {
				lgr.Printf("ERROR : st.SaveAllIntoFile returned the error %v\n", errSave)
			}
			lgr.Printf("saved %d metrics\n", count)
			st.ShowMetrics(ctx)
			if err = st.Close(); err != nil {
				lgr.Printf("ERROR : st.Close returned the error %v\n", err)
			}
		}()
		s = server.NewServer(st, cfg.Address, cfg.HashKey, cfg.CryptoKey)
	} else {
		defer func() {
			if err = db.ShowMetrics(ctx); err != nil {
				lgr.Printf("ERROR : db.ShowMetrics returned the error %v\n", err)
			}
			if err = db.Close(); err != nil {
				lgr.Printf("ERROR : db.Close returned the error %v\n", err)
			}
		}()
		storeInterval.Stop()
		s = server.NewServer(db, cfg.Address, cfg.HashKey, cfg.CryptoKey)
	}
	lgr.Printf("Start server")
	go func() {
		go lgr.Println(s.ListenAndServe())
	}()
	for {
		select {
		case <-storeInterval.C:
			count, errSave := st.SaveAllIntoFile()
			if errSave != nil {
				lgr.Printf("ERROR : st.SaveAllIntoFile returned the error %v\n", errSave)
			}
			lgr.Printf("saved %d metrics\n", count)
		case <-gracefulShutdown:
			lgr.Println("Server Get Signal!")
			storeInterval.Stop()
			if err = s.Shutdown(ctx); err != nil {
				lgr.Printf("ERROR : Server Shutdown error %v", err)
			} else {
				lgr.Printf("Server Shutdown ")
			}
			return
		}
	}
}
