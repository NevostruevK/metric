// Package commands обработка входных параметров для запуска программ
package commands

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
)

const maxRateLimit = 256

const (
	defAddress        = "127.0.0.1:8080"
	defStoreFile      = "/tmp/devops-metrics-db.json"
	defReportInterval = time.Second * 10
	defPollInterval   = time.Second * 2
	defStoreInterval  = time.Second * 300
	defRestore        = true
	defKey            = ""
	defDataBaseDSN    = ""
	defRateLimit      = 1
)

type Commands struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile      string        `env:"STORE_FILE" envDefault:""`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
	Key            string        `env:"KEY" envDefault:""`
	DataBaseDSN    string        `env:"DATABASE_DSN" envDefault:""`
	RateLimit      int           `env:"RATE_LIMIT" envDefault:"1"`
}

func GetAgentCommands() (*Commands, error) {
	addrPtr := flag.String("a", defAddress, "server address HOST:PORT")
	reportIntervalPtr := flag.Duration("r", defReportInterval, "report interval type : time.duration")
	pollIntervalPtr := flag.Duration("p", defPollInterval, "report interval type : time.duration")
	keyPtr := flag.String("k", defKey, "key for signing metrics")
	rateLimitPtr := flag.Int("l", defRateLimit, "requests count")
	flag.Parse()

	cmd := Commands{}
	err := env.Parse(&cmd)
	if _, ok := os.LookupEnv("ADDRESS"); !ok || err != nil {
		cmd.Address = *addrPtr
	}
	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok || err != nil {
		cmd.ReportInterval = *reportIntervalPtr
	}
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok || err != nil {
		cmd.PollInterval = *pollIntervalPtr
	}
	if _, ok := os.LookupEnv("KEY"); !ok || err != nil {
		cmd.Key = *keyPtr
	}
	if _, ok := os.LookupEnv("RATE_LIMIT"); !ok || err != nil {
		cmd.RateLimit = *rateLimitPtr
	}
	if cmd.RateLimit == 0 {
		cmd.RateLimit = 1
	}
	if cmd.RateLimit > maxRateLimit {
		cmd.RateLimit = maxRateLimit
	}
	return &cmd, err
}

func GetServerCommands() (*Commands, error) {
	addrPtr := flag.String("a", defAddress, "server address in format  Host:Port")
	restorePtr := flag.Bool("r", defRestore, "set if you need to load metric from file")
	storeIntervalPtr := flag.Duration("i", defStoreInterval, "store interval type : time.duration")
	storeFilePtr := flag.String("f", defStoreFile, "file for saving metrics")
	keyPtr := flag.String("k", defKey, "key for signing metrics")
	dataBasePtr := flag.String("d", defDataBaseDSN, "data base address")
	flag.Parse()

	cmd := Commands{}
	err := env.Parse(&cmd)

	if _, ok := os.LookupEnv("ADDRESS"); !ok || err != nil {
		cmd.Address = *addrPtr
	}
	if _, ok := os.LookupEnv("RESTORE"); !ok || err != nil {
		cmd.Restore = *restorePtr
	}
	if _, ok := os.LookupEnv("STORE_INTERVAL"); !ok || err != nil {
		cmd.StoreInterval = *storeIntervalPtr
	}
	if _, ok := os.LookupEnv("STORE_FILE"); !ok || err != nil {
		cmd.StoreFile = *storeFilePtr
	}
	if _, ok := os.LookupEnv("KEY"); !ok || err != nil {
		cmd.Key = *keyPtr
	}
	if _, ok := os.LookupEnv("DATABASE_DSN"); !ok || err != nil {
		cmd.DataBaseDSN = *dataBasePtr
	}
	return &cmd, err
}
