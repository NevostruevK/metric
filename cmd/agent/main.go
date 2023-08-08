package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/logger"
)

const shutDownTimeOut = time.Second * 3

var (
	buildVersion = "N/A"
	buildData    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)

	lgr.Println("Build version : " + buildVersion)
	lgr.Println("Build data    : " + buildData)
	lgr.Println("Build commit  : " + buildCommit)

	cfg := commands.GetAgentConfig()

	logger.LogCommands(cfg, false)

	complete := make(chan struct{})
	go client.StartAgent(ctx, cfg, complete)

	<-gracefulShutdown
	lgr.Println("Get Agent Signal!")
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), shutDownTimeOut)
	defer cancel()
	shutDownTimer := time.NewTimer(shutDownTimeOut)

	select {
	case <-shutDownTimer.C:
		lgr.Printf("shotdown with err %v", ctx.Err())
	case <-complete:
		lgr.Println("graceful shutdown")
	}
}
