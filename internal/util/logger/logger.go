// Package logger модуль логирования
package logger

import (
	"log"
	"os"

	"github.com/NevostruevK/metric/internal/util/commands"
)

var logWriter = os.Stdout

func LogCommands(cmd *commands.Config, isServer bool) {

	var logger *log.Logger
	if isServer {
		logger = log.New(logWriter, `server's flag : `, 0)
		logger.Printf("RESTORE = %t\n", cmd.Restore)
		logger.Printf("STORE_INTERVAL = %v\n", cmd.StoreInterval)
		logger.Printf("STORE_FILE = %s\n", cmd.StoreFile)
		logger.Printf("DATABASE_DSN = %s\n", cmd.DataBaseDSN)
		logger.Printf("TRUSTED_SUBNET = %s\n", cmd.TrustedSubnet)
	} else {
		logger = log.New(logWriter, `agent's flag : `, 0)
		logger.Printf("REPORT_INTERVAL = %v\n", cmd.ReportInterval)
		logger.Printf("POLL_INTERVAL = %v\n", cmd.PollInterval)
		logger.Printf("RATE_LIMIT = %v\n", cmd.RateLimit)
	}
	logger.Printf("ADDRESS = %s\n", cmd.Address)
	logger.Printf("KEY = %s\n", cmd.HashKey)
	logger.Printf("CRYPTO_KEY = %s\n", cmd.CryptoKey)
	logger.Printf("CERTIFICATE = %s\n", cmd.Certificate)
	logger.Printf("GRPC = %t\n", cmd.GRPC)
}

func NewLogger(name string, flags int) *log.Logger {
	return log.New(logWriter, name, flags)
}
