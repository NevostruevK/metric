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
	defCertificate    = ""
	defTrustedSubnet  = ""
	defCongig         = ""
	defRateLimit      = 1
	defRestore        = true
	defGRPC           = false
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
	usgCertificate    = "path to TLS certificate"
	usgTrustedSubnet  = "sub net in CIDR format"
	usgCongig         = "path to config file"
	usgRateLimit      = "requests count"
	usgRestore        = "restore value"
	usgGRPC           = "set true to use gRPC"
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
	flgCertificate    = "certificate"
	flgTrustedSubnet  = "t"
	flgConfig         = "config"
	flgConfigShort    = "c"
	flgRateLimit      = "l"
	flgRestore        = "r"
	flgGRPC           = "grpc"
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
	Certoificate   string        `env:"CERTIFICATE" envDefault:""`
	TrustedSubnet  string        `env:"TRUSTED_SUBNET" envDefault:""`
	RateLimit      int           `env:"RATE_LIMIT" envDefault:"1"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
	GRPC           bool          `env:"GRPC" envDefault:"false"`
}

func GetAgentConfig() *Config {
	var (
		reportInterval = flag.Duration(flgReportInterval, defReportInterval, usgReportInterval)
		pollInterval   = flag.Duration(flgPollInterval, defPollInterval, usgPollInterval)
		address        = flag.String(flgAddress, defAddress, usgAddress)
		hashKey        = flag.String(flgHashKey, defHashKey, usgHashKey)
		cryptoKey      = flag.String(flgCryptoKey, defCryptoKey, usgCryptoKey)
		certificate    = flag.String(flgCertificate, defCertificate, usgCertificate)
		rateLimit      = flag.Int(flgRateLimit, defRateLimit, usgRateLimit)
		gRPC           = flag.Bool(flgGRPC, defGRPC, usgGRPC)
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
		c.SetOption(WithAddress(value))
	}

	if value, ok := selectString("KEY", "", *hashKey); ok {
		c.SetOption(WithHashKey(value))
	}

	if value, ok := selectString("CRYPTO_KEY", "", *cryptoKey); ok {
		c.SetOption(WithCryptoKey(value))
	}

	if value, ok := selectString("CERTIFICATE", "", *certificate); ok {
		c.SetOption(WithCertificate(value))
	}

	if value, ok := selectDuration("REPORT_INTERVAL", time.Duration(0), e.ReportInterval, *reportInterval); ok {
		c.SetOption(WithReportInterval(duration.NewDuration(value)))
	}

	if value, ok := selectDuration("POOL_INTERVAL", time.Duration(0), e.PollInterval, *pollInterval); ok {
		c.SetOption(WithPollInterval(duration.NewDuration(value)))
	}

	if value, ok := selectInt("RATE_LIMIT", 0, e.RateLimit, *rateLimit); ok {
		c.SetOption(WithRateLimit(value))
	}

	if value, ok := selectBool("GRPC", defGRPC, e.GRPC, *gRPC); ok {
		c.SetOption(WithGRPC(value))
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
		certificate   = flag.String(flgCertificate, defCertificate, usgCertificate)
		trustedSubnet = flag.String(flgTrustedSubnet, defTrustedSubnet, usgTrustedSubnet)
		restore       = flag.Bool(flgRestore, defRestore, usgRestore)
		gRPC          = flag.Bool(flgGRPC, defGRPC, usgGRPC)
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
		c.SetOption(WithStoreInterval(duration.NewDuration(value)))
	}

	if value, ok := selectString("ADDRESS", "", *address); ok {
		c.SetOption(WithAddress(value))
	}

	if value, ok := selectString("STORE_FILE", "", *storeFile); ok {
		c.SetOption(WithStoreFile(value))
	}

	if value, ok := selectString("KEY", "", *hashKey); ok {
		c.SetOption(WithHashKey(value))
	}

	if value, ok := selectString("DATABASE_DSN", "", *dataBaseDSN); ok {
		c.SetOption(WithDataBaseDSN(value))
	}

	if value, ok := selectString("CRYPTO_KEY", "", *cryptoKey); ok {
		c.SetOption(WithCryptoKey(value))
	}

	if value, ok := selectString("CERTIFICATE", "", *certificate); ok {
		c.SetOption(WithCertificate(value))
	}

	if value, ok := selectString("TRUSTED_SUBNET", "", *trustedSubnet); ok {
		c.SetOption(WithTrustedSubnet(value))
	}

	if value, ok := selectBool("RESTORE", defRestore, e.Restore, *restore); ok {
		c.SetOption(WithRestore(value))
	}

	if value, ok := selectBool("GRPC", defGRPC, e.GRPC, *gRPC); ok {
		c.SetOption(WithGRPC(value))
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
