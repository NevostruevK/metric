// Package commands обработка входных параметров для запуска программ
package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NevostruevK/metric/internal/util/commands/duration"
	"github.com/caarlos0/env/v7"
)

const maxRateLimit = 256

const (
	defReportInterval = time.Second * 10
	defPollInterval   = time.Second * 2
	defStoreInterval  = time.Second * 300
	defAddress        = "127.0.0.1:8080"
	defStoreFile      = "/tmp/devops-metrics-db.json"
	defHashKey        = ""
	defDataBaseDSN    = ""
	defCryptoKey      = ""
	defCongig         = ""
	defRateLimit      = 1
	defRestore        = true
)

const (
	usgReportInterval = "report interval"
	usgPollInterval   = "poll interval"
	usgStoreInterval  = "store interval"
	usgAddress        = "server address HOST:PORT"
	usgStoreFile      = "store file"
	usgHashKey        = "key for signing metrics"
	usgDataBaseDSN    = "dsn"
	usgCryptoKey      = "path to private/public key"
	usgCongig         = "path to config file"
	usgRateLimit      = "requests count"
	usgRestore        = "restore value"
)

const (
	flgReportInterval = "r"
	flgPollInterval   = "p"
	flgStoreInterval  = "i"
	flgAddress        = "a"
	flgStoreFile      = "f"
	flgHashKey        = "k"
	flgDataBaseDSN    = "d"
	flgCryptoKey      = "crypto-key"
	flgConfig         = "config"
	flgConfigShort    = "c"
	flgRateLimit      = "l"
	flgRestore        = "r"
)

type Environment struct {
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile      string        `env:"STORE_FILE" envDefault:""`
	HashKey        string        `env:"KEY" envDefault:""`
	DataBaseDSN    string        `env:"DATABASE_DSN" envDefault:""`
	CryptoKey      string        `env:"CRYPTO_KEY" envDefault:""`
	//	Config         string        `env:"CONFIG" envDefault:""`
	RateLimit int  `env:"RATE_LIMIT" envDefault:"1"`
	Restore   bool `env:"RESTORE" envDefault:"true"`
}

func GetAgentConfig() *Config {
	var (
		reportInterval = flag.Duration(flgReportInterval, defReportInterval, usgReportInterval)
		pollInterval   = flag.Duration(flgPollInterval, defPollInterval, usgPollInterval)
		address        = flag.String(flgAddress, defAddress, usgAddress)
		hashKey        = flag.String(flgHashKey, defHashKey, usgHashKey)
		cryptoKey      = flag.String(flgCryptoKey, defCryptoKey, usgCryptoKey)
		rateLimit      = flag.Int(flgRateLimit, defRateLimit, usgRateLimit)
		config         = getFlagConfigValue()
	)
	flag.Parse()
	e := Environment{}
	err := env.Parse(&e)
	if err != nil {
		log.Fatal(err)
	}
	c := NewAgentConfig()
	fmt.Println(c)
	c.ReadConfig(getConfigPath("CONFIG", *config))
	fmt.Println(c)

	if value, ok := selectString("ADDRESS", "", *address); ok {
		c.setOption(withAddress(value))
	}

	if value, ok := selectString("KEY", "", *hashKey); ok {
		c.setOption(withHashKey(value))
	}

	if value, ok := selectString("CRYPTO_KEY", "", *cryptoKey); ok {
		c.setOption(withCryptoKey(value))
	}

	if value, ok := selectDuration("REPORT_INTERVAL", time.Duration(0), e.ReportInterval, *reportInterval); ok {
		c.setOption(withReportInterval(duration.NewDuration(value)))
	}

	if value, ok := selectDuration("POOL_INTERVAL", time.Duration(0), e.PollInterval, *pollInterval); ok {
		c.setOption(withPollInterval(duration.NewDuration(value)))
	}

	if value, ok := selectInt("RATE_LIMIT", 0, e.RateLimit, *rateLimit); ok {
		c.setOption(withRateLimit(value))
	}
	if c.RateLimit == 0 {
		c.RateLimit = 1
	}
	if c.RateLimit > maxRateLimit {
		c.RateLimit = maxRateLimit
	}
	return c
}

func GetServerConfig() *Config {
	var (
		storeInterval = flag.Duration(flgStoreInterval, defStoreInterval, usgStoreInterval)
		address       = flag.String(flgAddress, defAddress, usgAddress)
		storeFile     = flag.String(flgStoreFile, defStoreFile, usgStoreFile)
		hashKey       = flag.String(flgHashKey, defHashKey, usgHashKey)
		dataBaseDSN   = flag.String(flgDataBaseDSN, defDataBaseDSN, usgDataBaseDSN)
		cryptoKey     = flag.String(flgCryptoKey, defCryptoKey, usgCryptoKey)
		restore       = flag.Bool(flgRestore, defRestore, usgRestore)
		config        = getFlagConfigValue()
	)
	flag.Parse()
	e := Environment{}
	err := env.Parse(&e)
	if err != nil {
		log.Fatal(err)
	}

	c := NewServerConfig()
	c.ReadConfig(getConfigPath("CONFIG", *config))

	if value, ok := selectDuration("STORE_INTERVAL", time.Duration(0), e.StoreInterval, *storeInterval); ok {
		c.setOption(withStoreInterval(duration.NewDuration(value)))
	}

	if value, ok := selectString("ADDRESS", "", *address); ok {
		c.setOption(withAddress(value))
	}

	if value, ok := selectString("STORE_FILE", "", *storeFile); ok {
		c.setOption(withStoreFile(value))
	}

	if value, ok := selectString("KEY", "", *hashKey); ok {
		c.setOption(withHashKey(value))
	}

	if value, ok := selectString("DATABASE_DSN", "", *dataBaseDSN); ok {
		c.setOption(withDataBaseDSN(value))
	}

	if value, ok := selectString("CRYPTO_KEY", "", *cryptoKey); ok {
		c.setOption(withCryptoKey(value))
	}

	if value, ok := selectBool("RESTORE", defRestore, e.Restore, *restore); ok {
		c.setOption(withRestore(value))
	}
	return c
}

func getFlagConfigValue() *string {
	var config string
	flag.StringVar(&config, flgConfig, defCongig, usgCongig)
	flag.StringVar(&config, flgConfigShort, defCongig, usgCongig+" (shorthand)")
	return &config
}

func getConfigPath(env, flag string) string {
	if v, ok := os.LookupEnv(env); ok {
		return v
	}
	return flag
}

func selectString(env, def, flagString string) (string, bool) {
	if v, ok := os.LookupEnv(env); ok {
		return v, true
	}
	if flagString != def {
		return flagString, true
	}
	return "", false
}

func selectDuration(env string, def, envDuration, flgDuration time.Duration) (time.Duration, bool) {
	if _, ok := os.LookupEnv(env); ok {
		return envDuration, true
	}
	if flgDuration != def {
		return flgDuration, true
	}
	return time.Duration(0), false
}

func selectInt(env string, def, envInt, flgInt int) (int, bool) {
	if _, ok := os.LookupEnv(env); ok {
		return envInt, true
	}
	if flgInt != def {
		return flgInt, true
	}
	return 0, false
}

func selectBool(env string, def, envInt, flgBool bool) (bool, bool) {
	if _, ok := os.LookupEnv(env); ok {
		return envInt, true
	}
	if flgBool != def {
		return flgBool, true
	}
	return false, false
}
