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

const shotDownTimeOut = time.Second * 3

var (
	buildVersion = "N/A"
	buildData    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	complete := make(chan struct{})
	ctx := context.Background()

	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	lgr.Println("Build version : " + buildVersion)
	lgr.Println("Build data    : " + buildData)
	lgr.Println("Build commit  : " + buildCommit)

	lgr.Println(`Get server's flags`)

	cfg := commands.GetServerConfig()
	logger.LogCommands(cfg, true)

	storeInterval := time.NewTicker(cfg.StoreInterval.Duration)
	st := &storage.MemStorage{}

	s := &http.Server{}

	lgr.Println(`Init database`)
	db, err := db.NewDB(ctx, cfg.DataBaseDSN)
	if err != nil || cfg.DataBaseDSN == "" {
		lgr.Println(`Init Memory storage`)
		st = storage.NewMemStorage(cfg.Restore, cfg.StoreInterval.Duration == 0, cfg.StoreFile)
		defer func() {
			if err = st.Close(ctx); err != nil {
				lgr.Printf("ERROR : st.Close returned the error %v\n", err)
			}
		}()
		s = server.NewServer(st, cfg.Address, cfg.HashKey, cfg.CryptoKey)
	} else {
		defer func() {
			if err = db.Close(ctx); err != nil {
				lgr.Printf("ERROR : db.Close returned the error %v\n", err)
			}
		}()
		storeInterval.Stop()
		s = server.NewServer(db, cfg.Address, cfg.HashKey, cfg.CryptoKey)
	}

	lgr.Printf("Start server")
	go func() {
		lgr.Println(s.ListenAndServe())
		if db.Init {
			err = db.Close(ctx)
		} else {
			err = st.Close(ctx)
		}
		if err != nil {
			lgr.Printf("ERROR : storage.Close returned the error %v\n", err)
		}
		complete <- struct{}{}
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
			shutdownCtx, cancel := context.WithTimeout(context.Background(), shotDownTimeOut)
			defer cancel()
			if err = s.Shutdown(shutdownCtx); err != nil {
				lgr.Printf("ERROR : Server Shutdown error %v", err)
			} else {
				lgr.Printf("Server Shutdown ")
			}
			select {
			case <-shutdownCtx.Done():
				lgr.Printf("shotdown with err %v", shutdownCtx.Err())
			case <-complete:
				lgr.Println("graceful shutdown")
			}
			return
		}
	}
}
