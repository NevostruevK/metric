package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/db"
	grpcserver "github.com/NevostruevK/metric/internal/grpc/server"
	restserver "github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/logger"
)

const shutDownTimeOut = time.Second * 3

var (
	buildVersion = "N/A"
	buildData    = "N/A"
	buildCommit  = "N/A"
)

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	complete := make(chan struct{})
	ctx := context.Background()

	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	lgr.Println("Build version : " + buildVersion)
	lgr.Println("Build data    : " + buildData)
	lgr.Println("Build commit  : " + buildCommit)

	cfg := commands.GetServerConfig()
	logger.LogCommands(cfg, true)

	storeInterval := time.NewTicker(cfg.StoreInterval.Duration)

	var st storage.Repository
	var s server

	lgr.Println(`Init database`)
	st, err := db.NewDB(ctx, cfg.DataBaseDSN)
	if err != nil || cfg.DataBaseDSN == "" {
		lgr.Println(`Init Memory storage`)
		st = storage.NewMemStorage(cfg.Restore, cfg.StoreInterval.Duration == 0, cfg.StoreFile)
		defer func() {
			if err = st.Close(ctx); err != nil {
				lgr.Printf("ERROR : st.Close returned the error %v\n", err)
			}
		}()
	} else {
		defer func() {
			if err = st.Close(ctx); err != nil {
				lgr.Printf("ERROR : db.Close returned the error %v\n", err)
			}
		}()
		storeInterval.Stop()
	}

	if cfg.GRPC {
		s, err = grpcserver.NewServer(st, cfg)
	} else {
		s, err = restserver.NewServer(st, cfg)
	}
	if err != nil {
		lgr.Fatalln(err)
	}

	lgr.Printf("Start server")
	go func() {
		lgr.Println(s.ListenAndServe())
		err = st.Close(ctx)
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
			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutDownTimeOut)
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
