package logger

import (
	"log"
	"os"

	"github.com/NevostruevK/metric/internal/util/commands"
)

var logWriter = os.Stdout
var Server = &log.Logger{}

func LogCommands(cmd *commands.Commands, isServer bool, err error) {

	var logger *log.Logger
	if isServer {
		logger = log.New(logWriter, `server's flag : `, 0)
		logger.Printf("RESTORE = %t\n", cmd.Restore)
		logger.Printf("STORE_INTERVAL = %v\n", cmd.StoreInterval)
		logger.Printf("STORE_FILE = %s\n", cmd.StoreFile)
		logger.Printf("DATABASE_DSN = %s\n", cmd.DataBaseDSN)
	} else {
		logger = log.New(logWriter, `agent's flag : `, 0)
		logger.Printf("REPORT_INTERVAL = %v\n", cmd.ReportInterval)
		logger.Printf("POLL_INTERVAL = %v\n", cmd.PollInterval)
	}
	logger.Printf("ADDRESS = %s\n", cmd.Address)
	logger.Printf("KEY = %s\n", cmd.Key)
	if err != nil {
		logger.Printf("ERROR : read environment with the error: %v\n", err)
	}
}

func NewLogger(name string, flags int) *log.Logger {
	return log.New(logWriter, name, flags)
}

func NewServer(name string, flags int) *log.Logger {
	Server = log.New(logWriter, name, flags)
	return Server
}
