package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/logger"
)

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

	cmd, err := commands.GetAgentCommands()
	logger.LogCommands(cmd, false, err)

	go client.StartAgent(ctx, cmd)

	select {
	case <-gracefulShutdown:
		lgr.Println("Get Agent Signal!")
		cancel()
	case <-ctx.Done():
	}
}
