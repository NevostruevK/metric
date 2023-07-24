package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/metric/internal/client"
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
	ctx, cancel := context.WithCancel(context.Background())

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	lgr := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)

	lgr.Println("Build version : " + buildVersion)
	lgr.Println("Build data    : " + buildData)
	lgr.Println("Build commit  : " + buildCommit)

	lgr.Println(`Get server's flags`)

	cfg := commands.GetAgentConfig()
	fmt.Println(cfg)

	logger.LogCommands(cfg, false)

	complete := make(chan struct{})
	go client.StartAgent(ctx, cfg, complete)

	<-gracefulShutdown
	lgr.Println("Get Agent Signal!")
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), shotDownTimeOut)
	defer cancel()

	select {
	case <-ctx.Done():
		lgr.Printf("shotdown with err %v", ctx.Err())
		return
	case <-complete:
		lgr.Println("graceful shutdown")
		return
	}
}
