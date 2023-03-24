package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
)
const(
	defAddress 			= "127.0.0.1:8080"
	defStoreFile 		= "/tmp/devops-metrics-db.json"
	defReportInterval 	= time.Second*10
	defPollInterval 	= time.Second*2
	defStoreInterval 	= time.Second*300
	defRestore			= true
)

type Commands struct {
	Address        	string       	`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile	 	string  		`env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	ReportInterval 	time.Duration 	`env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   	time.Duration 	`env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval	time.Duration  	`env:"STORE_INTERVAL" envDefault:"300s"`
	Restore			bool			`env:"RESTORE" envDefault:"true"`
	parsingError   	bool
}

func GetAgentCommands() *Commands{
	cmd := Commands{parsingError: false}
	addrPtr := flag.String("a", defAddress, "server address HOST:PORT")
	reportIntervalPtr := flag.Duration("r", defReportInterval, "report interval type : time.duration")
	pollIntervalPtr := flag.Duration("p", defPollInterval, "report interval type : time.duration")
	flag.Parse()

	if err := env.Parse(&cmd); err!=nil{
		fmt.Printf("Agent read environment with the error: %+v\n", err)
		cmd.parsingError = true
	}
	if _, err := os.LookupEnv("ADDRESS"); !err || cmd.parsingError{
		cmd.Address = *addrPtr
	}
	if _, err := os.LookupEnv("REPORT_INTERVAL"); !err || cmd.parsingError{
		cmd.ReportInterval = *reportIntervalPtr
	}
	if _, err := os.LookupEnv("POLL_INTERVAL"); !err || cmd.parsingError{
		cmd.PollInterval = *pollIntervalPtr
	}
	return &cmd
}

func GetServerCommands() *Commands{
	cmd := Commands{parsingError: false}
	addrPtr := flag.String("a", defAddress, "server address in format  Host:Port")
	restorePtr := flag.Bool("r", defRestore, "set if you need to load metric from file")
	storeIntervalPtr := flag.Duration("i", defStoreInterval, "store interval type : time.duration")
	storeFilePtr := flag.String("f", defStoreFile, "file for saving metrics")
	flag.Parse()

	if err := env.Parse(&cmd); err!=nil{
		fmt.Printf("Agent read environment with the error: %+v\n", err)
		cmd.parsingError = true
	}
	if _, err := os.LookupEnv("ADDRESS"); !err || cmd.parsingError{
		cmd.Address = *addrPtr
	}
	if _, err := os.LookupEnv("RESTORE"); !err || cmd.parsingError{
		cmd.Restore = *restorePtr
	}
	if _, err := os.LookupEnv("STORE_INTERVAL"); !err || cmd.parsingError{
		cmd.StoreInterval = *storeIntervalPtr
	}
	if _, err := os.LookupEnv("STORE_FILE"); !err || cmd.parsingError{
		cmd.StoreFile = *storeFilePtr
	}
	return &cmd
}